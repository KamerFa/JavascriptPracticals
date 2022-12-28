const productsJSON = new Object(
        {
            "products": [

            ]
        }
    );
const filters = getFilters();
const productsContainer = document.getElementById('products');
const filtersContainer = document.getElementById('filters');
const filteredProducts = new Object(
        {
        "products": [

            ]
        }
    );

    function getFilters(){
        const items = [];
        const filter = "filter";

        for(let [key, value] of Object.entries(localStorage)){
            if(key.startsWith(filter)){
                items.push(value);
            }
        } 
        
        console.log("FILTERED FROM LOCAL", items);
        return items;
    }


    function filterProducts(obj, filters, filteredProducts) {
        let products = obj.products;

        for (let i = 0; i < products.length; i++) {
            const product = products[i];
            if(filters.includes(product.product_type)){
                filteredProducts.products.push(product);
            }
        } 
    }


    async function getProducts(obj)
        {
        const response = await fetch('./products.json');
        const data = await response.json();
        
        let products = data.products;


        for (let i = 0; i < products.length; i++) {
            const product = products[i];
            obj.products.push(product);
        }
    }

    function renderCards(obj){

            const productsMap = obj.products.map(product => {
                {
                    const card = document.createElement('div');
                    card.className = ("product-card col-4");
                    card.innerHTML =
                        `<div>
                            <img src='${product.images[0].src}' class="img-fluid">
                            ${product.title}
                        </div>`;

                    return card;
                }
            });
            productsMap.forEach(card => productsContainer.appendChild(card));
        };
        
    

    function renderFilters(obj){

        const uniqueTypes = new Object(
            {
                "types":[
                    {
                        "index": 1,
                        "id": "Dresses",
                    },
                    {   
                        "index": 2,
                        "id": "Jacket",
                    },
                    {
                        "index": 3,
                        "id": "Trousers",
                    },
                    {
                        "index": 4,
                        "id": "Blouses",
                    },
                    {   
                        "index": 5,
                        "id": "T-Shirt",
                    },
                    {
                        "index": 6,
                        "id": "Denim",
                    }
                ]
            }

        );
        console.log(uniqueTypes);


        // const getUniqueProductTypes = function (obj) {

        //     for (const product of obj.products) {
        //         uniqueTypes.add(product.product_type);
        //     }
        //     console.log(uniqueTypes, typeof (uniqueTypes));
        // }

        // getUniqueProductTypes(obj);

        const filtersMap = uniqueTypes.types.map(filterU => {
            {
                const filter = document.createElement('div');
                filter.className = ("filter border border-danger");
                filter.addEventListener('click', function(event){
                    localStorage.setItem("filter"+`${filterU.index}`,`${filterU.id}`);
                    console.log(filters);
                });


                filter.innerHTML =
                    `<button value=${filterU.id}>
                        ${filterU.id}
                    </button>`;

                return filter;
            }
        });
        console.log(filtersMap);
        filtersMap.forEach(filter => filtersContainer.appendChild(filter));
    };


async function App(fetched, filters, filtered) {




    await getProducts(fetched);
    filterProducts(fetched, filters, filtered);

    renderFilters(fetched);

    
    if(filters.length != 0){
        renderCards(filtered);
    }else{
        renderCards(fetched);
    }
    


}



App(productsJSON, filters, filteredProducts);



