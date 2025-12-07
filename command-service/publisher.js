// const amqp = require('amqplib/callback_api')

// amqp.connect('amqp://guest:guest@localhost:5672/',(err,connection)=>{
//     if(err) throw err
//     connection.createChannel((err,channel)=>{
//         if(err) throw err
//         let queueName = 'create_product'
//         let message = 'this is a create product'
//         channel.assertQueue(queueName,{
//             durable:false
//         })
//         channel.sendToQueue(queueName,Buffer.from(message))
//         setTimeout(()=>{
//             connection.close()
//         },1000)
//     })
// })