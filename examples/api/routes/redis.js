const express = require("express")
const router = express.Router()

const redisn = require("redis");
var redis = require('promise-redis')();
const client = redis.createClient("redis://redis");


client.on("error", function(error) {
  console.error(error);
});

/* Bulk add to cache in 1s */
router.get('/addBulk/:requests', function(req, res, next) {
    const { requests=100 } = req.params;
    
    const times = []
    
    const interval = requests / 1000
  
    let gotRequests = 0
  
    for ( let i = 0; i < requests; i += 1 ) {
  
      let end
      setTimeout(() => {
        let start = new Date()
        client.set("abc", "value")
        .then((response) => {
            end = new Date() 
            total = end.getTime() - start.getTime()
            times.push(total)
            console.log(total)
            gotRequests += 1
            if ( gotRequests == requests ) {
                let avgTime = 0;
                for ( let i = 0; i < times.length; i ++ ) {
                    avgTime += times[i]
                }
                avgTime = avgTime / times.length
                res.status(200).send(`time: ${times} 
                avg: ${ avgTime}`)
            }
        })
      }, interval  )
    }
  });

router.get('/store/:key', async (req, res) => {
    const { key } = req.params;
    const value = req.query;

    client.set(key, "value")
    .then(() => {
        return res.send('Success');
    })

});

router.get('/:key', async (req, res) => {
    const { key } = req.params;
    client.get(key)
    .then((response) => {
        return res.json(response);
    })
});


module.exports = router 