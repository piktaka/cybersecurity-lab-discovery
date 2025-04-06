const express = require("express");
const router = express.Router();
const bcrypt = require("bcryptjs");
const User = require("../models/User");
const fs = require("fs").promises; // Use promises for cleaner async operations
const path = require("path");

// Signup route
router.post("/signup", async (req, res) => {
  try {
    const { username, password, email } = req.body;

    if (!username || !password || password.length >= 12) {
      return res.status(400).json({
        success: false,
        message:
          "Username and password are required or need to be 12 characters",
      });
    }

    const existingUser = await User.findOne({
      $or: [{ username: username }, { email: email }],
    });
    if (existingUser) {
      return res.status(400).json({
        success: false,
        message: "Username or email already exists",
      });
    }

    const salt = await bcrypt.genSalt(10);
    const hashedPassword = await bcrypt.hash(password, salt);

    const user = await User.create({
      username,
      email,
      password: hashedPassword,
      lastPasswordChange: new Date(),
      passwordHistory: [{ password: hashedPassword, changedAt: new Date() }],
    });

    // Create the file after successful user creation
    const signupFilePath = path.join("/opt", "validate1");
    const signupFileContent = "challenge3";

    try {
      await fs.writeFile(signupFilePath, signupFileContent, "utf8");
      console.log(`Signup file created successfully at ${signupFilePath}`);
    } catch (signupFileError) {
      console.error("Error creating signup file:", signupFileError);
    }

    res.status(201).json({
      success: true,
      message: "User created successfully",
    });
  } catch (error) {
    console.error("Signup error:", error);
    res.status(500).json({
      success: false,
      message: error.message || "Error creating user",
    });
  }
});

// Login route
router.post("/login", async (req, res) => {
  try {
    const { username, password } = req.body;

    if (!username || !password) {
      return res.status(400).json({
        success: false,
        message: "Username and password are required",
      });
    }

    const user = await User.findOne({
      $or: [{ username: username }, { email: username }],
    });
    if (!user) {
      return res.status(401).json({
        success: false,
        message: "Invalid credentials",
      });
    }

    const isMatch = await bcrypt.compare(password, user.password);
    if (!isMatch) {
      return res.status(401).json({
        success: false,
        message: "Invalid credentials",
      });
    }

    if (user.isPasswordExpired()) {
      return res.status(403).json({
        success: false,
        message: "Your password has expired. Please change your password.",
        requiresPasswordChange: true,
      });
    }

    res.json({
      success: true,
      message: "Login successful",
      user: {
        username: user.username,
        email: user.email,
        role: user.role,
        lastPasswordChange: user.lastPasswordChange,
        createdAt: user.createdAt,
      },
    });
  } catch (error) {
    console.error("Login error:", error);
    res.status(500).json({
      success: false,
      message: "Error during login",
    });
  }
});

// Get all users (admin route)
router.get("/users", async (req, res) => {
  try {
    // Get all users excluding sensitive information
    const users = await User.find(
      {},
      {
        password: 0,
        passwordHistory: 0,
      }
    );

    // Calculate expiry time and minutes left for each user
    const usersWithExpiry = users.map((user) => {
      const expiryTime = new Date(user.lastPasswordChange);
      expiryTime.setMinutes(
        expiryTime.getMinutes() + user.passwordExpiryMinutes
      );

      const now = new Date();
      const minutesLeft = Math.floor((expiryTime - now) / (1000 * 60));

      return {
        ...user.toObject(),
        expiryTime,
        minutesLeft,
      };
    });

    res.json({
      success: true,
      users: usersWithExpiry,
    });
  } catch (error) {
    console.error("Error fetching users:", error);
    res.status(500).json({
      success: false,
      message: "Error fetching users",
    });
  }
});

// Change password route
router.post("/change-password", async (req, res) => {
  try {
    const { username, currentPassword, newPassword } = req.body;

    if (!username || !currentPassword || !newPassword) {
      return res.status(400).json({
        success: false,
        message: "All fields are required",
      });
    }

    const user = await User.findOne({
      $or: [{ username: username }, { email: username }],
    });
    if (!user) {
      return res.status(401).json({
        success: false,
        message: "User not found",
      });
    }

    const isMatch = await bcrypt.compare(currentPassword, user.password);
    if (!isMatch) {
      return res.status(401).json({
        success: false,
        message: "Current password is incorrect",
      });
    }

    const isPasswordReused = await user.isPasswordPreviouslyUsed(newPassword);
    if (isPasswordReused) {
      return res.status(400).json({
        success: false,
        message:
          "This password was previously used. Please choose a different password.",
      });
    }

    const salt = await bcrypt.genSalt(10);
    const hashedPassword = await bcrypt.hash(newPassword, salt);

    user.passwordHistory.push({
      password: user.password,
      changedAt: user.lastPasswordChange,
    });
    user.password = hashedPassword;
    user.lastPasswordChange = new Date();
    await user.save();

    // Create the file after successful password change
    const changePasswordFilePath = path.join("/opt", "validate2");
    const changePasswordFileContent = "challenge3";

    try {
      await fs.writeFile(
        changePasswordFilePath,
        changePasswordFileContent,
        "utf8"
      );
      console.log(
        `Password change file created successfully at ${changePasswordFilePath}`
      );
    } catch (changePasswordFileError) {
      console.error(
        "Error creating password change file:",
        changePasswordFileError
      );
    }

    res.json({
      success: true,
      message: "Password changed successfully",
    });
  } catch (error) {
    res.status(500).json({
      success: false,
      message: error.message || "Error changing password",
    });
  }
});

module.exports = router;
