const express = require('express')
const db = require('./db')
const app = express()
const fs = require('fs');
const path = require('path');
const api = require('./routes/index')
const { connectRabbitMQ } = require('./services/rabbitMqService')

app.get('/initiate', async (req, res) => {
  try {
    const sqlPath = path.join(__dirname, '..', 'docker', 'postgres', 'init.sql');
    const sql = fs.readFileSync(sqlPath, 'utf8');
    console.log('Running SQL script... FOR tables initialization');
    const result = await db.query(sql);
    console.log(result)
  } catch (err) {
    console.error(err);
    res.status(500).send('Internal Server Error');
  }
});

app.use(express.urlencoded({ extended: false })); //built in middlware to handlw urlencoded data (form data)
app.use(express.json());

// app.get('/',async (req,res)=>{
//     try {
//     const result = await db.query('SELECT * FROM product');
//     res.json(result.rows);
//   } catch (err) {
//     console.error(err);
//     res.status(500).send('Internal Server Error');
//   }  

// })
app.use('/api', api);

connectRabbitMQ().then(() => {
  app.listen(5000, () => {
    console.log("hi")
  })
});
