const express = require('express')
const router = express.Router()
const clientController = require('../controllers/clientController')

router.post('/',clientController.create)
// router.put('/:id',clientController.update)
router.delete('/:id',clientController.dlt)
// router.get('/',clientController.getAll)

module.exports = router;