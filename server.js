
// File: server.js
const express = require('express');
const mongoose = require('mongoose');
const cors = require('cors');
const path = require('path');
const app = express();

// Middleware
app.use(cors());
app.use(express.json());

// MongoDB Connection
mongoose.connect('mongodb://mongo:27017/comments_db', {
    useNewUrlParser: true,
    useUnifiedTopology: true,
    serverSelectionTimeoutMS: 30000

});

// Comment Schema
const commentSchema = new mongoose.Schema({
    text: { type: String, required: true },
    author: { type: String, required: true },
    userId: { type: String, required: true },
    timestamp: { type: Date, default: Date.now }
});

const Comment = mongoose.model('Comment', commentSchema);

app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, 'public', 'index.html'));
});
// Routes
app.get('/api/comments', async (req, res) => {
    try {
        const comments = await Comment.find().sort({ timestamp: -1 });
        res.json(comments);
    } catch (error) {
        res.status(500).json({ error: 'Error fetching comments' });
    }
});

app.post('/api/comments', async (req, res) => {
    try {
        const { text, author, userId } = req.body;
        const comment = new Comment({ text, author, userId });
        await comment.save();
        res.status(201).json(comment);
    } catch (error) {
        res.status(500).json({ error: 'Error creating comment' });
    }
});

const PORT ="5000" ///  process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});