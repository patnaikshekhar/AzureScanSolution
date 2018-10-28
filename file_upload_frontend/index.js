const express = require('express')
const azure = require('azure-storage')
const app = express()
const path = require('path')
const QUARANTINE_CONTAINER = 'quarantine'
const CLEAN_CONTAINER = 'clean'
const VIRUS_CONTAINER = 'virus'

var multer  = require('multer')
var upload = multer({ dest: 'uploads/' })

app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, 'index.html'))
})

app.get('/fileStatus', async (req, res) => {
    try {
        var blobService = azure.createBlobService()
        res.json({
            status: 'success',
            quarantine: await getBlobsForContainer(blobService, QUARANTINE_CONTAINER),
            clean: await getBlobsForContainer(blobService, CLEAN_CONTAINER),
            virus: await getBlobsForContainer(blobService, VIRUS_CONTAINER)
        })          
    } catch(e) {
        console.error(e)
        res.json({
            status: 'error'
        })
    }
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
            res.redirect('/')
        }
    })
})

app.listen(process.env.PORT || 80, () => {
    console.log(`AZURE_STORAGE_ACCOUNT=${process.env.AZURE_STORAGE_ACCOUNT} AZURE_STORAGE_ACCESS_KEY=${process.env.AZURE_STORAGE_ACCESS_KEY}`)
    console.log('Server listening')
})

function getBlobsForContainer(blobService, containerName) {
    return new Promise((resolve, reject) => {
        blobService.listBlobsSegmented(containerName, null, (err, result) => {
            if (err) {
                reject(err)
            } else {
                resolve(result.entries)
            }
        })
    })
}