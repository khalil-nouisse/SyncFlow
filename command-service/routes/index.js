const express = require('express')
const router = express.Router()

const productRoute = require('./product')
const clientRoute = require('./client')
const commandeRoute = require('./commande')
// const commandeProdRoute = require('./commande-prod')

const swaggerUi = require('swagger-ui-express');
const swaggerSpecs = require('../swagger');

router.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerSpecs));
router.use('/product', productRoute);
router.use('/client',clientRoute);
router.use('/commande', commandeRoute);
// router.use('/commandeProd',commandeProdRoute);

module.exports = router