const amqp = require('amqplib');

let channel = null;
let connection = null;

const connectRabbitMQ = async () => {
    try {
        // Connect to the RabbitMQ server
        connection = await amqp.connect('amqp://guest:guest@localhost:5672/');
        
        // Create a channel
        channel = await connection.createChannel();
        
        console.log("✅ RabbitMQ Connected and Channel Created");
        
        // Handle connection close/error events to auto-reconnect (Optional but recommended)
        connection.on('error', (err) => {
            console.error('RabbitMQ connection error', err);
            // logic to reconnect could go here
        });

    } catch (error) {
        console.error("Failed to connect to RabbitMQ:", error);
        process.exit(1); // Exit if MQ is essential for startup
    }
};

const publishToQueue = async (queueName, data) => {
    if (!channel) {
        console.error("RabbitMQ channel not initialized. Call connectRabbitMQ first.");
        return;
    }

    try {
        // Ensure queue exists (idempotent)
        await channel.assertQueue(queueName, { durable: false });
        
        // Convert object to Buffer and send
        const messageBuffer = Buffer.from(JSON.stringify(data));
        channel.sendToQueue(queueName, messageBuffer);
        
        console.log(`[x] Sent to ${queueName}:`, data);
    } catch (error) {
        console.error("Error publishing message:", error);
    }
};

module.exports = {
    connectRabbitMQ,
    publishToQueue
};