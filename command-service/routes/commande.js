const express = require('express')
const router = express.Router()
const commandeController = require('../controllers/commandeController')

router.post('/', commandeController.create)
router.put('/:id', commandeController.update)
router.delete('/:id', commandeController.deleteOrder)
router.get('/', commandeController.getAll)

module.exports = router;