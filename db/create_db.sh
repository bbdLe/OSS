#!/bin/bash
curl localhost:9200/metadata -XPUT -d '{"mappings":{"objects":{"properties":{"name":{"type":"string","index":"not analyzed"},"version":{"type":"integer"},"size":{"type":"integer"},"hash":{"type":"string"}}}}}'  -H "Content-Type: application/json"
