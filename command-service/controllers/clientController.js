const db = require('../db')
const { publishToQueue } = require('../services/rabbitMqService');
const getAll = async (req,res)=>{
    try{
    const result = await db.query('SELECT * FROM product')
    console.log(result.rows)
    }catch(err){
        console.error(err);
    }
}

const create =  async (req,res)=>{
    try{
        const {first_name,last_name,address} = req.body;
        const result = await db.query('INSERT INTO CLIENT (firstname,lastname,address) values($1,$2,$3) RETURNING *',[first_name,last_name,address])
        const newClient = result.rows[0];

        const event = {
            event_type: "CLIENT_CREATED",
            payload: newClient
        };
        await publishToQueue('product_events', event);
        res.status(201).json({ message: "Client created", client: newClient });
    }
    catch(err){
        console.error(err);
        res.status(500).json({ error: "Internal Server Error" });
    }

}

const dlt = async (req,res)=>{
    try{
        id_client = req.params.id
        const result = await db.query('DELETE FROM CLIENT where id_client = $1 RETURNING *',[id_client])
        const toDeleteClient = result.rows[0]

        const event = {
            event_type: "CLIENT_DELETED",
            payload: toDeleteClient
        };
        await publishToQueue('product_events', event);
        res.status(201).json({ message: "Client deleted", client: toDeleteClient });
    }catch(err){
        console.error(err)
    }
}
module.exports = {getAll,create,dlt}