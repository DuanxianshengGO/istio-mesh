#!/bin/bash

# Test script for Istio API endpoints
# This script tests the basic functionality of the Istio API

BASE_URL="http://localhost:8080/api/v1/istio"
CLUSTER_NAME="local-cluster"
NAMESPACE="default"

echo "Testing Istio API endpoints..."

# Test VirtualService endpoints
echo "1. Testing VirtualService endpoints..."

# List VirtualServices
echo "  - Listing VirtualServices..."
curl -s -X GET "${BASE_URL}/${CLUSTER_NAME}/virtualservices" | jq '.' > /dev/null
if [ $? -eq 0 ]; then
    echo "    ✓ List VirtualServices - OK"
else
    echo "    ✗ List VirtualServices - FAILED"
fi

# List VirtualServices with namespace
echo "  - Listing VirtualServices in namespace ${NAMESPACE}..."
curl -s -X GET "${BASE_URL}/${CLUSTER_NAME}/namespaces/${NAMESPACE}/virtualservices" | jq '.' > /dev/null
if [ $? -eq 0 ]; then
    echo "    ✓ List VirtualServices with namespace - OK"
else
    echo "    ✗ List VirtualServices with namespace - FAILED"
fi

# Test DestinationRule endpoints
echo "2. Testing DestinationRule endpoints..."

# List DestinationRules
echo "  - Listing DestinationRules..."
curl -s -X GET "${BASE_URL}/${CLUSTER_NAME}/destinationrules" | jq '.' > /dev/null
if [ $? -eq 0 ]; then
    echo "    ✓ List DestinationRules - OK"
else
    echo "    ✗ List DestinationRules - FAILED"
fi

# List DestinationRules with namespace
echo "  - Listing DestinationRules in namespace ${NAMESPACE}..."
curl -s -X GET "${BASE_URL}/${CLUSTER_NAME}/namespaces/${NAMESPACE}/destinationrules" | jq '.' > /dev/null
if [ $? -eq 0 ]; then
    echo "    ✓ List DestinationRules with namespace - OK"
else
    echo "    ✗ List DestinationRules with namespace - FAILED"
fi

# Test Gateway endpoints
echo "3. Testing Gateway endpoints..."

# List Gateways
echo "  - Listing Gateways..."
curl -s -X GET "${BASE_URL}/${CLUSTER_NAME}/gateways" | jq '.' > /dev/null
if [ $? -eq 0 ]; then
    echo "    ✓ List Gateways - OK"
else
    echo "    ✗ List Gateways - FAILED"
fi

# List Gateways with namespace
echo "  - Listing Gateways in namespace ${NAMESPACE}..."
curl -s -X GET "${BASE_URL}/${CLUSTER_NAME}/namespaces/${NAMESPACE}/gateways" | jq '.' > /dev/null
if [ $? -eq 0 ]; then
    echo "    ✓ List Gateways with namespace - OK"
else
    echo "    ✗ List Gateways with namespace - FAILED"
fi

# Test Traffic Analytics endpoint
echo "4. Testing Traffic Analytics endpoint..."
curl -s -X GET "${BASE_URL}/${CLUSTER_NAME}/traffic-analytics?namespace=${NAMESPACE}" | jq '.' > /dev/null
if [ $? -eq 0 ]; then
    echo "    ✓ Traffic Analytics - OK"
else
    echo "    ✗ Traffic Analytics - FAILED"
fi

echo "API testing completed!"

# Test creating a sample VirtualService
echo "5. Testing VirtualService creation..."

SAMPLE_VS='{
  "apiVersion": "networking.istio.io/v1beta1",
  "kind": "VirtualService",
  "metadata": {
    "name": "test-vs",
    "namespace": "default"
  },
  "spec": {
    "hosts": ["test-service"],
    "http": [
      {
        "route": [
          {
            "destination": {
              "host": "test-service",
              "subset": "v1"
            }
          }
        ]
      }
    ]
  }
}'

echo "  - Creating test VirtualService..."
curl -s -X POST "${BASE_URL}/${CLUSTER_NAME}/namespaces/${NAMESPACE}/virtualservices" \
     -H "Content-Type: application/json" \
     -d "${SAMPLE_VS}" > /dev/null
if [ $? -eq 0 ]; then
    echo "    ✓ Create VirtualService - OK"
    
    # Try to get the created VirtualService
    echo "  - Getting created VirtualService..."
    curl -s -X GET "${BASE_URL}/${CLUSTER_NAME}/namespaces/${NAMESPACE}/virtualservices/test-vs" | jq '.' > /dev/null
    if [ $? -eq 0 ]; then
        echo "    ✓ Get VirtualService - OK"
        
        # Clean up - delete the test VirtualService
        echo "  - Cleaning up test VirtualService..."
        curl -s -X DELETE "${BASE_URL}/${CLUSTER_NAME}/namespaces/${NAMESPACE}/virtualservices/test-vs" > /dev/null
        if [ $? -eq 0 ]; then
            echo "    ✓ Delete VirtualService - OK"
        else
            echo "    ✗ Delete VirtualService - FAILED"
        fi
    else
        echo "    ✗ Get VirtualService - FAILED"
    fi
else
    echo "    ✗ Create VirtualService - FAILED"
fi

echo "All tests completed!"
