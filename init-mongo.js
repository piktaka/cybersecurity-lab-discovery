db = db.getSiblingDB('comments_db'); // Switch to `comments_db`

db.createCollection('comments'); // Create the `comments` collection


// Insert pre-populated comments
db.comments.insertMany([
    {
        text: "The 5G lab environment is incredibly useful for testing our network configurations. Really impressed with the setup!",
        author: "Sarah Chen",
        userId: "user1",
        timestamp: new Date("2024-01-01")
    },
    {
        text: "OpenStack integration works flawlessly. Saved us weeks of infrastructure setup time.",
        author: "Mike Johnson",
        userId: "user2",
        timestamp: new Date("2024-01-01")
    },
    {
        text: "Thanks for the feedback! We're continuously working on improving the lab environments.",
        author: "Brown Smith",
        userId: "user3",
        timestamp: new Date("2024-01-01")
    },
    {
        text: "Great platform for hands-on practice. The 4G to 5G migration scenarios are particularly helpful.",
        author: "Elena Rodriguez",
        userId: "user4",
        timestamp: new Date("2024-01-01")
    }
]);
