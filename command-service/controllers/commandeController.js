const db = require('../db')
const { publishToQueue } = require('../services/rabbitMqService');

/**
 * @swagger
 * /api/commande:
 *   post:
 *     summary: Create a new order
 *     tags: [Orders]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - id_client
 *             properties:
 *               id_client:
 *                 type: string
 *     responses:
 *       201:
 *         description: Order created
 *       500:
 *         description: Internal Server Error
 */
const create = async (req, res) => {
    try {
        const { id_client,products } = req.body;
        // Basic validation
        if (!id_client || !products) {
            return res.status(400).json({ error: "id_client and products are both required" });
        }

        // Insert into Postgres
        // Returning * to get the generated id_commande and date_commande
        const result = await db.query(
            'INSERT INTO commande (id_client) VALUES ($1) RETURNING *',
            [id_client]
        );
        const newOrder = result.rows[0];
        products.forEach(async (product)=>{
            await db.query(
                'INSERT INTO commande_prod (id_commande,id_produit,prix_unitaire,quantity) VALUES($1,$2,$3,$4)',
                [newOrder.id_commande,product.id_produit,product.prix_unitaire,product.quantity]
            );
        })


        const payload = {
            id_client,
            products,
            newOrder
        }
        // Publish Event
        const event = {
            event_type: "ORDER_CREATED",
            payload: payload
        };
        await publishToQueue('product_events', event);

        res.status(201).json({ message: "Order created", order: newOrder });
    } catch (err) {
        console.error(err);
        res.status(500).json({ error: "Internal Server Error" });
    }
}

/**
 * @swagger
 * /api/commande/{id}:
 *   put:
 *     summary: Update an order
 *     tags: [Orders]
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             properties:
 *               date_commande:
 *                 type: string
 *     responses:
 *       200:
 *         description: Order updated
 *       404:
 *         description: Order not found
 */
const update = async (req, res) => {
    try {
        const { id } = req.params; // id_commande
        const { date_commande } = req.body; // allow updating date manually for demo

        const result = await db.query(
            'UPDATE commande SET date_commande = $1 WHERE id_commande = $2 RETURNING *',
            [date_commande || new Date(), id]
        );

        if (result.rows.length === 0) {
            return res.status(404).json({ error: "Order not found" });
        }

        const updatedOrder = result.rows[0];

        // Publish Event
        const event = {
            event_type: "ORDER_UPDATED",
            payload: updatedOrder
        };
        await publishToQueue('product_events', event);

        res.json({ message: "Order updated", order: updatedOrder });
    } catch (err) {
        console.error(err);
        res.status(500).json({ error: "Internal Server Error" });
    }
}

/**
 * @swagger
 * /api/commande/{id}:
 *   delete:
 *     summary: Delete an order
 *     tags: [Orders]
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: string
 *     responses:
 *       200:
 *         description: Order deleted
 *       404:
 *         description: Order not found
 */
const deleteOrder = async (req, res) => {
    try {
        const { id } = req.params;

        const result = await db.query(
            'DELETE FROM commande WHERE id_commande = $1 RETURNING *',
            [id]
        );

        if (result.rows.length === 0) {
            return res.status(404).json({ error: "Order not found" });
        }

        const deletedOrder = result.rows[0];

        // Publish Event
        const event = {
            event_type: "ORDER_DELETED",
            payload: deletedOrder
        };
        await publishToQueue('product_events', event);

        res.json({ message: "Order deleted", order: deletedOrder });

    } catch (err) {
        console.error(err);
        res.status(500).json({ error: "Internal Server Error" });
    }
}

/**
 * @swagger
 * /api/commande:
 *   get:
 *     summary: Get all orders (Postgres View - usually this endpoint might not be needed if Query service exists, but exists for admin?)
 *     tags: [Orders]
 *     responses:
 *       200:
 *         description: List of orders
 */
const getAll = async (req, res) => {
    try {
        const result = await db.query('SELECT * FROM commande');
        res.json(result.rows);
    } catch (err) {
        console.error(err);
        res.status(500).json({ error: "Internal Server Error" });
    }
}

module.exports = { create, update, deleteOrder, getAll }