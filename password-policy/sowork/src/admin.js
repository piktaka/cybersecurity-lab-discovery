// src/initAdmin.js
const bcrypt = require('bcryptjs');
const User = require('./models/User');

async function createDefaultAdmin() {
    try {
        // Admin credentials
        const adminEmail = 'admin@sowork.com'
        const adminUsername = 'admin';
        const adminPassword = 'admin123';
        
        // Check if admin already exists
        const existingAdmin = await User.findOne({
              $or: [
                { username: adminUsername },
                { email: adminEmail }
              ]
            });
        if (existingAdmin) {
            console.log('Admin account already exists');
            return;
        }

        // Hash the password
        const salt = await bcrypt.genSalt(10);
        const hashedPassword = await bcrypt.hash(adminPassword, salt);

        // Create admin user
        await User.create({
            username: adminUsername,
            email:adminEmail,
            password: hashedPassword,
            role: 'admin',
            lastPasswordChange: new Date(),
            passwordHistory: [{ password: hashedPassword, changedAt: new Date() }]
        });

        console.log('Default admin account created successfully');
    } catch (error) {
        console.error('Error creating admin account:', error);
    }
}

module.exports = createDefaultAdmin;
