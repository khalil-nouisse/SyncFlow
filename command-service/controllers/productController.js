const db = require('../db')
const { publishToQueue } = require('../services/rabbitMqService');

/**
 * @swagger
 * /api/product:
 *   post:
 *     summary: Create a new product
 *     tags: [Products]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - p_desc
 *               - qte
 *             properties:
 *               p_desc:
 *                 type: string
 *               qte:
 *                 type: number
 *     responses:
 *       201:
 *         description: Product created successfully
 *       500:
 *         description: Internal Server Error
 */
const create = async (req, res) => {
    try {
        const { p_desc, qte } = req.body;
        const productQte = parseFloat(qte);
        console.log(productQte)

        // 1. Create the Event Wrapper
        const event = {
            event_type: "PRODUCT_CREATED",
            payload: {
                p_desc: p_desc,
                qte: productQte
            }
        };

        // 2. Publish to 'product_events' queue
        await publishToQueue('product_events', event);

        const result = db.query('INSERT INTO PRODUCT (p_description,stock_quantity) values($1,$2)', [p_desc, productQte])
        console.log(result.rows)
        res.status(201).json({ message: "Product created and event published", product: { p_desc, qte: productQte } });
    }
    catch (err) {
        console.error(err)
        res.status(500).json({ error: "Internal Server Error" });
    }

}
const getAll = async (req, res) => {
    try {
        const result = await db.query('SELECT * FROM product')
        console.log(result.rows)
        res.json(result.rows)
    } catch (err) {
        console.error(err);
        res.status(500).json({ error: "Internal Server Error" });
    }
}

module.exports = { create, getAll }