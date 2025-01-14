const mongoose = require('mongoose');

const userSchema = new mongoose.Schema({
  username: {
    type: String,
    required: [true, 'Username is required'],
    unique: true,
    trim: true,
    minlength: [3, 'Username must be at least 3 characters long']
  },
  email: {
    type: String,
    required: [true, 'Email is required'],
    unique: true,
    trim: true,
    match: [/\S+@\S+\.\S+/, 'Please provide a valid email address']
  },
  password: {
    type: String,
    required: [true, 'Password is required']
  },
  role: {
    type: String,
    enum: ['user', 'admin'],
    default: 'user'
  },
  passwordHistory: [{
    password: String,
    changedAt: {
      type: Date,
      default: Date.now
    }
  }],
  lastPasswordChange: {
    type: Date,
    default: Date.now
  },
  passwordExpiryMinutes: {
    type: Number,
    default: 5 // 5 minutes for testing
  },
  createdAt: {
    type: Date,
    default: Date.now
  }
});


// Method to check if password is expired
userSchema.methods.isPasswordExpired = function() {
  const now = new Date();
  const expiryTime = new Date(this.lastPasswordChange);
  expiryTime.setMinutes(expiryTime.getMinutes() + this.passwordExpiryMinutes);
  return now > expiryTime;
};

// Method to check if password was previously used
userSchema.methods.isPasswordPreviouslyUsed = async function(newPassword) {
  const bcrypt = require('bcryptjs');
  
  // Check current password
  const isSameAsCurrent = await bcrypt.compare(newPassword, this.password);
  if (isSameAsCurrent) return true;

  // Check password history
  for (const historyEntry of this.passwordHistory) {
    const isMatch = await bcrypt.compare(newPassword, historyEntry.password);
    if (isMatch) return true;
  }
  
  return false;
};

module.exports = mongoose.model('User', userSchema);
