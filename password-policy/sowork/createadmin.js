// scripts/createAdmin.js
require('dotenv').config();
const mongoose = require('mongoose');
const bcrypt = require('bcryptjs');
const User = require('./src/models/User');

async function createAdmin() {
    try {
        // Connect to MongoDB
        await mongoose.connect(process.env.MONGODB_URI, {
            useNewUrlParser: true,
            useUnifiedTopology: true
        });

        // Admin credentials
        const adminUsername = 'admin@sowork.com';
        const adminPassword = 'admin123';

        // Check if admin exists
        const existingAdmin = await User.findOne({ username: adminUsername });
        
        if (existingAdmin) {
            console.log('Admin account already exists');
            process.exit(0);
        }

        // Hash password
        const salt = await bcrypt.genSalt(10);
        const hashedPassword = await bcrypt.hash(adminPassword, salt);

        // Create admin
        await User.create({
            username: adminUsername,
            password: hashedPassword,
            role: 'admin',
            lastPasswordChange: new Date(),
            passwordHistory: [{ password: hashedPassword, changedAt: new Date() }]
        });

        console.log('Admin account created successfully');
    } catch (error) {
        console.error('Error:', error);
    } finally {
        await mongoose.connection.close();
        process.exit(0);
    }
}

createAdmin();
