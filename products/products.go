package main
import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
	//"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	//"github.com/hyperledger/fabric/core/chaincode/shim"
	//pb "github.com/hyperledger/fabric/protos/peer"
	
)

// PropertyDetails smart contract to show the product details
type ProductDetailsSmartContract struct {
	contractapi.Contract
	}


	// Product  Obejct
	
	type Product struct {
		ID	string `json:"id"`
		Name	string `json:"name"`
		Value	int `json:"value"`
	}
	
	// This function helps to Add new Product
func (pc *ProductDetailsSmartContract) AddProduct(ctx contractapi.TransactionContextInterface, id string,  name string,  value int) error {
    productJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return fmt.Errorf("Failed to read the data from world state", err)
    }
	
    if productJSON != nil {
		return fmt.Errorf("the product %s already exists", id)
    }
	
	prod := Product{
		ID:            id,
		Name:          name,
		Value: 		   value,
	}
	
	productBytes, err := json.Marshal(prod)	
	if err != nil {
		return err
	}

    return ctx.GetStub().PutState(id, productBytes)
}


// This function returns all the existing products 
func (pc *ProductDetailsSmartContract) QueryAllProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
	productIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer productIterator.Close()

	var products []*Product
	for productIterator.HasNext() {
		productResponse, err := productIterator.Next()
		if err != nil {
			return nil, err
		}

		var product *Product
		err = json.Unmarshal(productResponse.Value, &product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// This function helps to query the product by Id
func (pc *ProductDetailsSmartContract) QueryProductById(ctx contractapi.TransactionContextInterface, id string) (*Product, error) {
    productJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("Failed to read the data from world state", err)
    }
	
    if productJSON == nil {
		return nil, fmt.Errorf("the property %s does not exist", id)
    }
	
	var product *Product
	err = json.Unmarshal(productJSON, &product)
	
	if err != nil {
		return nil, err
	}
	return product, nil
}


func (pc *ProductDetailsSmartContract) DeleteProductById(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	if len(id) == 0 {
		return "", fmt.Errorf("Please provide correct product Id")
	}

	return ctx.GetStub().GetTxID(), ctx.GetStub().DelState(id)
}

func (pc *ProductDetailsSmartContract) DeleteAllProducts(ctx contractapi.TransactionContextInterface)  error {
	var p *Product

	p = new (Product)
	fmt.Println("Records are deleted")
	fmt.Println(p)
}

/*
function to restrict the access 
func (s *SmartContract) restictedMethod(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	

	val, ok, err := cid.GetAttributeValue(APIstub, "role")
	if err != nil {
		// There was an error trying to retrieve the attribute
		shim.Error("Error while retriving attributes")
	}
	if !ok {
		
		shim.Error("Client identity doesnot posses the attribute")
	}

	if val != "approver" {
		fmt.Println("Attribute role: " + val)
		return shim.Error("Only user with role as APPROVER have access this method!")
	} else {
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1")
		}

		productAsBytes, _ := APIstub.GetState(args[0])
		return shim.Success(productAsBytes)
	}

}
 */

func main() {
    productDetailsSmartContract := new(ProductDetailsSmartContract)

    cc, err := contractapi.NewChaincode(productDetailsSmartContract)

    if err != nil {
        panic(err.Error())
    }

    if err := cc.Start(); err != nil {
        panic(err.Error())
    }
}
