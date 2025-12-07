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
        const {p_desc,qte} = req.body;
        const productQte = parseFloat(qte);
        console.log(productQte)
        await publishToQueue('create_product',req.body)

       const result = db.query('INSERT INTO PRODUCT (p_description,stock_quantity) values($1,$2)',[p_desc,productQte])
       console.log(result.rows)
    }
    catch(err){
        console.error(err)
    }

}
module.exports = {getAll,create}