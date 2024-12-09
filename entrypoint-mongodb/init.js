const fs = require('fs');

db = db.getSiblingDB('mytheresadb');
db.Product.createIndex({ category: 1 });
db.Product.createIndex({ price: 1 });
db.Product.createIndex({ category: 1, price: 1 });

const fileProducts = JSON.parse(fs.readFileSync('/docker-entrypoint-initdb.d/products.json', 'utf8'));
const products = fileProducts.products;

const batchSize = 10000;
for (let i = 0; i < products.length; i += batchSize) {
    const batch = products.slice(i, i + batchSize);
    db.Product.insertMany(batch);
    print(`Insert batch ${Math.ceil(i / batchSize) + 1} of ${Math.ceil(products.length / batchSize)}`);
}

db.Product.updateOne(
    { sku: "000003" },
    { $set: { sku_discount: 15 } }
);

db.DiscountCategory.createIndex({ category: 1 });

const fileDiscounts = JSON.parse(fs.readFileSync('/docker-entrypoint-initdb.d/discounts.json', 'utf8'));
const discounts = fileDiscounts.discounts_category;

db.DiscountCategory.insertMany(discounts);

print("Data succesfully load on mongodb.");
