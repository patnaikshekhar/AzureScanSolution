const express = require('express')
const azure = require('azure-storage')
const app = express()
const path = require('path')
const QUARANTINE_CONTAINER = 'quarantine'

var multer  = require('multer')
var upload = multer({ dest: 'uploads/' })

app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, 'index.html'))
})

app.post('/upload', upload.single('file'), (req, res) => {
    var blobService = azure.createBlobService()
    blobService.createBlockBlobFromLocalFile(QUARANTINE_CONTAINER, req.file.originalname, req.file.path, function(error, result, response) {
        if (error) {
            res.json({
                status: 'error',
                error: error
            })
        } else {
            res.json({
                status: 'success'
            })
        }
    })
})

app.listen(process.env.PORT || 80, () => {
    console.log(`AZURE_STORAGE_ACCOUNT=${process.env.AZURE_STORAGE_ACCOUNT} AZURE_STORAGE_ACCESS_KEY=${process.env.AZURE_STORAGE_ACCESS_KEY}`)
    console.log('Server listening')
})