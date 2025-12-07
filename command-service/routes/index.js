const express = require('express')
const router = express.Router()

const productRoute = require('./product')
// const clientRoute = require('./client')
// const commandeRoute = require('./commande')
// const commandeProdRoute = require('./commande-prod')

router.use('/product',productRoute);
// router.use('/client',clientRoute);
// router.use('/commande',commandeRoute);
// router.use('/commandeProd',commandeProdRoute);

module.exports = router