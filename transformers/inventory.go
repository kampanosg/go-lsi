package transformers

import (
    "github.com/kampanosg/go-lsi/types"
)

func FromCategoryDbRowToDomain(id, squareId, name string, version int64) types.Category {
    return types.Category {
        Id: id,
        SquareId: squareId,
        Name: name,
        Version: version,
    }
}

func FromProductDbRowToDomain(id, squareId, squareVarId, categoryId, squareCategoryId, title, barcode, sku string, price float64, version int64) types.Product {
    return types.Product {
        Id: id,
        SquareId: squareId,
        SquareVarId: squareVarId, 
        CategoryId: categoryId,
        SquareCategoryId: squareCategoryId,
        Title: title,
        Barcode: barcode,
        SKU: sku,
        Price: price,
        Version: version,
    }
}
