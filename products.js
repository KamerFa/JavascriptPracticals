const productsContainer = document.getElementById('products');
const filtersContainer = document.getElementById('filters');
const filteredProducts = new Object(
        {
        "products": [

            ]
        }
    );
const productsJSON = new Object(
    {
        "products": [

        ]
    }
);
    async function getProducts(obj) {
        const response = await fetch('./products.json');
        const data = await response.json();

        let products = data.products;


        for (let i = 0; i < products.length; i++) {
            const product = products[i];
            obj.products.push(product);
        }
    }

    function getOptions(obj){

        let products = obj.products;

        for(let i = 0; i < products.length; i++){
            let product = products[i];
            
                
            const card = document.createElement('div');
            card.id = 'card';
            card.className='col-4';

            const image = document.createElement('div');
            image.innerHTML = `
                <img src=${product.images[0].src} alt="" srcset="" class="img-fluid">
            `;

            const title = document.createElement('div');
            title.className = "product-title d-flex";
            title.innerHTML =
            `<h5 class="fs-5">${product.title}<h5>`;
            
            const colors = document.createElement('div');
            colors.id = 'product-color-options'
            colors.className = "d-flex flex-wrap";

            const sizes = document.createElement('div');
            sizes.id = 'product-size-options';
            sizes.className = "d-flex"

            productsContainer.appendChild(card);
            
            card.appendChild(image);
            card.appendChild(title);
            card.appendChild(colors);
            card.appendChild(sizes);

            for(const option of product.options){
                let values = option.values;
                let j = 0;

                switch (option.name) {
                    case 'Color':
                        values.forEach(value => {
                            // Create a new element for each value
                            const element = document.createElement('div');
                            element.id = option.name;
                            element.className = `${value} mx-1`;
                            element.innerHTML = `
                                    ${value}
                            `;
                            colors.appendChild(element);
                        });
                        break;
                    case 'Size':
                        values.forEach(value => {
                            // Create a new element for each value
                            const element = document.createElement('div');
                            element.id = option.name;
                            element.className = `${value} mx-1`
                            element.innerHTML = `
                                <div class="${value}">
                                    ${value}
                                </div>
                            `;
                            sizes.appendChild(element);
                        });
                        break;
                    default:
                        values.forEach(value => {
                            // Create a new element for each value
                            const element = document.createElement('div');
                            element.id = option.name;
                            element.innerHTML = `
                                <div class="${value}">
                                    ${value}
                                </div>
                            `;
                            card.appendChild(element);
                        });
                }
            }
            };
    };



    // function filterProducts(obj, filters, filteredProducts) {
    //     let products = obj.products;

    //     for (let i = 0; i < products.length; i++) {
    //         const product = products[i];
    //         if(filters.includes(product.product_type)){
    //             filteredProducts.products.push(product);
    //         }
    //     } 
    // }




    function renderCards(obj){

            const productsMap = obj.products.map(product => {
                {
                    const card = document.createElement('div');
                    card.className = ("product-card col-4");
                    card.innerHTML =
                        `<div>
                            <img src=${product.images[0].src} class="img-fluid">
                            ${product.title}
                        </div>`;

                    return card;
                }
            });
            productsMap.forEach(card => productsContainer.appendChild(card));
        };
    


async function App(fetched, filters, filtered) {

    await getProducts(fetched);

    const options = new Array();

    getOptions(fetched, options);
    // filterProducts(fetched, filters, filtered);


    // if(filters.length != 0){
    //     renderCards(filtered);
    // }else{
    //     renderCards(fetched);
    // }

};



App(productsJSON, filters, filteredProducts);



