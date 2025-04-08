// File: server.js
const express = require("express");
const mongoose = require("mongoose");
const cors = require("cors");
const path = require("path");
const app = express();

// Middleware
app.use(cors());
app.use(express.json());

// MongoDB Connection
mongoose.connect("mongodb://mongo:27017/comments_db", {
  useNewUrlParser: true,
  useUnifiedTopology: true,
  serverSelectionTimeoutMS: 30000,
});

// Comment Schema
const commentSchema = new mongoose.Schema({
  text: { type: String, required: true },
  author: { type: String, required: true },
  userId: { type: String, required: true },
  timestamp: { type: Date, default: Date.now },
});

const Comment = mongoose.model("Comment", commentSchema);

app.get("/", (req, res) => {
  res.sendFile(path.join(__dirname, "public", "index.html"));
});
// Routes
app.get("/api/comments", async (req, res) => {
  try {
    const comments = await Comment.find().sort({ timestamp: -1 });
    res.json(comments);
  } catch (error) {
    res.status(500).json({ error: "Error fetching comments" });
  }
});

const XSS_PATTERNS = [
  "<script",
  "</script",
  "<iframe",
  "</iframe",
  "<img",
  "<svg",
  "onerror=",
  "javascript:",
  "vbscript:",
  "data:",
  "eval(",
  "document.",
  "window.",
  "alert(",
  "prompt(",
  "confirm(",
  "onload=",
  "onmouseover=",
  "onclick=",
  // Add more patterns as needed
];

app.post("/api/comments", async (req, res) => {
  try {
    const { text, author, userId } = req.body;
    let isXSSLikely = false;
    for (const pattern of XSS_PATTERNS) {
      if (text.toLowerCase().includes(pattern.toLowerCase())) {
        isXSSLikely = true;
        break;
      }
    }
    if (isXSSLikely) {
      const validationFile = "/opt/validation/validation1";
      const fileContent = "challenge4";

      try {
        // Check if the file exists
        try {
          await fs.access(validationFile);
          console.log(`Validation file already exists at ${validationFile}`);
        } catch (error) {
          // File does not exist, so create it
          await fs.writeFile(validationFile, fileContent, "utf8");
          console.log(`Validation file created at ${validationFile}`);
        }
      } catch (fileError) {
        console.error("Error creating/accessing validation file:", fileError);
        // Handle the error appropriately
      }
    }

    const comment = new Comment({ text, author, userId });
    await comment.save();
    res.status(201).json(comment);
  } catch (error) {
    res.status(500).json({ error: "Error creating comment" });
  }
});

const PORT = "5000"; ///  process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
