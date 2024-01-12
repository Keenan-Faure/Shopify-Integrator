import {useEffect} from 'react';
import {useState} from "react";
import $ from 'jquery';
import '../../CSS/login.css';

import FroalaEditorComponent from 'react-froala-wysiwyg';
import 'froala-editor/css/froala_style.min.css';
import 'froala-editor/css/froala_editor.pkgd.min.css';

function Add_Product()
{
    const[inputs, setInputs] = useState({});

    const handleChange = (event) =>
    {
        const name = event.target.name;
        const value = event.target.value;
        setInputs(values => ({...values, [name]: value}))
    }
    const AddProduct = (event) =>
    {
        event.preventDefault();

        let element = document.querySelector(".fr-element").firstChild;
        console.log(element.innerHTML);

        let Object = 
        {
            product_code: inputs.product_code, 
            title: inputs.product_title, 
            body_html: element.innerHTML, 
            category: inputs.product_category, 
            vendor: inputs.product_vendor,
            product_type: inputs.product_type, 
            variants: 
            [
                {
                    sku: inputs.variant_sku, 
                    option1: inputs.variant_option1, 
                    option2: inputs.variant_option2, 
                    option3: inputs.variant_option3, 
                    barcode: inputs.variant_barcode,
                    variant_price_tiers: 
                    [
                        {
                            name: inputs.price_tier_name,
                            value: inputs.price_tier_value
                        },
                        {
                            name: inputs.compareprice_tier_name,
                            value: inputs.price_tier_value2
                        }
                    ],
                    variant_quantities: 
                    [
                        {
                            name: inputs.quantity_warehouse_name,
                            value: parseInt(inputs.warehouse_quantity)
                        }
                    ]
                }
            ],
            options: 
            [
                {
                    value: inputs.product_options
                },
                {
                    value: inputs.product_options2
                }
            ]
        };

        console.log(Object);
        const api_key = localStorage.getItem('api_key');
        
        $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
        $.post("http://localhost:8080/api/products", JSON.stringify(Object),[], 'json')
        .done(function( _data) 
        {
            console.log(_data);
        })
        .fail( function(xhr) 
        {
            alert(xhr.responseText);
        });
        
    }

    useEffect(() =>
    {
        window.onload = function(event)
        {
            let navbar = document.getElementById("navbar");
            navbar.style.display = "none";
        }
        /* Fix any incorrect elements */
        let navigation = document.getElementById("navbar");
        let modal = document.getElementById("model");
        modal.style.display = "block";
        navigation.style.animation = "MoveRight 1.2s ease";
        navigation.style.position = "fixed";
        navigation.style.left = "0%";
        navigation.style.width = "100%";

        function openPage(pageName) 
        {
            var i, tabcontent, tablinks;
            tabcontent = document.getElementsByClassName("tabcontent");
            for (i = 0; i < tabcontent.length; i++) 
            {
                tabcontent[i].style.display = "none";
            }
            tablinks = document.getElementsByClassName("tablink");
            for (i = 0; i < tablinks.length; i++) 
            {
                tablinks[i].style.backgroundColor = "";
            }

            document.getElementById("_" + pageName).style.display = "block";
            document.getElementById(pageName).style.backgroundColor = "rgb(72, 101, 128)";
            document.getElementById(pageName).style.color = "black";
            
        }

        let home = document.getElementById("Product");
        home.addEventListener("click", () =>
        {
            openPage('Product');
        });

        let defaul = document.getElementById("Variants");
        defaul.addEventListener("click", () =>
        {
            openPage('Variants');
        });

        document.getElementById("Product").click();

        /* When the user clicks on the return button */
        let close = document.querySelector(".rtn-button");
        let navbar = document.getElementById("navbar");
        let details = document.getElementById("detailss");
        navbar.style.display = "none";
        close.addEventListener("click", ()=> 
        {
            close.style.display = "none";
            details.style.animation = "Fadeout 0.5s ease-out";
            
            window.location.href = "/products";
        });

        const input = document.getElementById("image_Input");
        const output = document.getElementById("image_Output");
        let label = document.getElementById("label");
        let del = document.getElementById("delete");
        let imagesArray = []
        del.style.display = "none";

        input.addEventListener("change", () => 
        {
            const file = input.files
            imagesArray.push(file[0]);
            displayImages();
            label.style.display = "none";
            del.style.display = "block";
        });

        function displayImages() 
        {
            let images = "";
            imagesArray.forEach((image, index) => 
            {
              images += `<div class="image" id = "imgg">
                          <img class = "detailed-image" src="${URL.createObjectURL(image)}">
                        </div>`
            })
            output.innerHTML = images;
        }

        del.addEventListener("click", () =>
        {
            document.getElementById("imgg").remove();
            label.style.display = "";
            del.style.display = "none";
        });

        /* Styles for fr-box element */
        setTimeout(() => 
        {
            let element_box = document.querySelector(".fr-box");
            element_box.style.width = "95%"; element_box.style.left = "50%"; element_box.style.transform = "translate(-50%)";
        }, 20);
        


    }, []);

    return (
        <>
            <div className = 'modal1' id = "model" style={{zIndex: '2'}}>
                <form className = 'modal-content' style ={{opacity: '1'}} method = 'post' onSubmit={(event) => AddProduct(event)} autoComplete='off' id = 'form1' encType="multipart/form-data" noValidate>
                <div id = "detailss">
                    <div className = 'rtn-button'></div>
                    <div className = "button-holder" style = {{position: 'absolute', width: '71%', zIndex:'3', left:'29%', backgroundColor: ''}}>
                        <button type = "button" className="tablink" id = "Product" style ={{left: '-14%', width:'95px'}}>Product</button>
                        <button type = "button" className="tablink" id ="Variants" style ={{left: '-14%', width:'95px'}}>Variants</button>
                    </div>
                
                    <div className="tabcontent" id="_Product" >
                        <div className = "details-details">
                            <div className = "detailed" style = {{backgroundColor: 'transparent'}}>
                                <div className = "details-title">
                                    <input type = '_text' style ={{fontSize:'20px', width: '500px'}} placeholder = "Product Title" name = "product_title" value = {inputs.product_title || ""}  
                                    onChange = {handleChange} required></input>
                                </div>
                                <table>
                                    <tbody>
                                        <tr>
                                            <th>Product Category</th>
                                            <th>Product Code</th>
                                        </tr>
                                        <tr>
                                            <td><input type = '_text' style = {{width: '300px'}} placeholder = "Product Category" name = "product_category" 
                                            value = {inputs.product_category || ""} onChange = {handleChange} required></input></td>
                                            <td><input type = '_text' style = {{width: '300px'}} placeholder = "Product Code" name = "product_code" 
                                            value = {inputs.product_code || ""} onChange = {handleChange} required></input></td>
                                        </tr>
                                        <tr>
                                            <th>Product Type</th>
                                            <th>Product Vendor</th>
                                        </tr>
                                        <tr>
                                            <td><input type = '_text' style = {{width: '300px'}} placeholder = "Product Type" name = "product_type" 
                                            value = {inputs.product_type || ""} onChange = {handleChange} required></input></td>
                                            <td><input type = '_text' style = {{width: '300px'}} placeholder = "Product Vendor" name = "product_vendor" 
                                            value = {inputs.product_vendor || ""} onChange = {handleChange} required></input></td>
                                        </tr>
                                    </tbody>
                                </table> 
                                <div className = "details-description">Product Description</div>
                                    <FroalaEditorComponent name = "product_description" tag = 'textarea' 
                                    placeholderText = "Product Description" value = {inputs.product_description || ""}/>

                                    
                            </div>
                            <div className = "details-left" style = {{backgroundColor: 'transparent'}}/>
                            <div className = "details-right" style = {{backgroundColor: 'transparent'}}>
                                <input id = "image_Input" name = "product_image" value = {inputs.product_image || ""} onChange = {handleChange}
                                style = {{position: 'relative', top: '5px'}} type="file" accept="image/jpeg, image/png, image/jpg" hidden/>
                                <label id = "label" htmlFor = "image_Input">Upload Product Image</label>
                                <span id = "delete" >&times;</span>
                                <output id = "image_Output"></output>
                            </div>
                        </div>
                    </div>

                    <div className="tabcontent" id="_Variants" >
                        <div className = "details-details">
                            <div className = "detailed" style = {{backgroundColor: 'transparent'}}>
                                <div className = "details-title">Variants</div>
                                <div className = "variants" id="_variants" >
                                <table>
                                        <tbody>
                                            <tr>
                                                <th>Variant Barcode</th>
                                                <th>Variant ID</th>
                                                <th>Variant SKU</th>
                                            </tr>
                                            <tr>
                                                <td><input type = '_text' style = {{width: '100%'}} placeholder = "Variant Barcode" name = "variant_barcode" 
                                                value = {inputs.variant_barcode || ""} onChange = {handleChange} required></input></td>
                                                <td><input type = '_text' style = {{width: '100%'}} placeholder = "Variant ID" name = "variant_id" 
                                                value = {inputs.variant_id || ""} onChange = {handleChange} required></input></td>
                                                <td><input type = '_text' style = {{width: '100%'}} placeholder = "Variant SKU" name = "variant_sku" 
                                                value = {inputs.variant_sku || ""} onChange = {handleChange} required></input></td>
                                            </tr>
                                        </tbody>
                                    </table>
                                    <table>
                                        <tbody>
                                            <tr>
                                                <th>Option 1</th>
                                                <th>Option 2</th>
                                                <th>Option 3</th>
                                            </tr>
                                            <tr>
                                                <td><input type = '_text' style = {{width: '100%'}} placeholder = "Variant Option 1" name = "variant_option1" 
                                                value = {inputs.variant_option1 || ""} onChange = {handleChange} required></input></td>
                                                <td><input type = '_text' style = {{width: '100%'}} placeholder = "Variant Option 2" name = "variant_option2" 
                                                value = {inputs.variant_option2 || ""} onChange = {handleChange} ></input></td>
                                                <td><input type = '_text' style = {{width: '100%'}} placeholder = "Variant Option 3" name = "variant_option3" 
                                                value = {inputs.variant_option3 || ""} onChange = {handleChange} ></input></td>
                                            </tr>
                                        </tbody>
                                    </table>

                                    <div className = "details-description">Variant Price Tiers</div>
                                        <table>
                                            <tbody>
                                                <tr>
                                                    <th>Name</th>
                                                    <th>Value</th>
                                                </tr>
                                                <tr>
                                                    <td><input type = '_text' style = {{width: '100%'}} placeholder = "Selling Price" name = "price_tier_name" 
                                                    value = {inputs.price_tier_name || ""} onChange = {handleChange} required></input></td>
                                                    <td><input type = '_text' style = {{width: '100%'}} placeholder = "Price Tier Price" name = "price_tier_value" 
                                                    value = {inputs.price_tier_value || ""} onChange = {handleChange} ></input></td>
                                                </tr>
                                                <tr>
                                                    <td><input type = '_text' style = {{width: '100%'}} placeholder = "Compare-to-Price" name = "_compareprice_tier_name" 
                                                    value = {inputs._compareprice_tier_name || ""} onChange = {handleChange} required></input></td>
                                                    <td><input type = '_text' style = {{width: '100%'}} placeholder = "Price Tier Price" name = "price_tier_value2" 
                                                    value = {inputs.price_tier_value2 || ""} onChange = {handleChange} ></input></td>
                                                </tr>
                                            </tbody>
                                        </table>

                                    <div className = "details-description">Variant Quantites</div>
                                        <table>
                                            <tbody>
                                                <tr>
                                                    <th>Warehouse Name</th>
                                                    <th>Quantity</th>
                                                </tr>
                                                <tr>
                                                    <td><input type = '_text' style = {{width: '100%'}} placeholder = "Warehouse Name" name = "quantity_warehouse_name" 
                                                    value = {inputs.quantity_warehouse_name || ""} onChange = {handleChange} required></input></td>
                                                    <td><input type = '_text' style = {{width: '100%'}} placeholder = "Warehouse Quantity" name = "warehouse_quantity" 
                                                    value = {inputs.warehouse_quantity || ""} onChange = {handleChange} ></input></td>
                                                </tr>
                                                <tr>
                                                    <td><input type = '_text' style = {{width: '100%'}} placeholder = "Warehouse 2 Name" name = "quantity_warehouse_name2" 
                                                    value = {inputs.quantity_warehouse_name2 || ""} onChange = {handleChange} required></input></td>
                                                    <td><input type = '_text' style = {{width: '100%'}} placeholder = "Warehouse 2 Quantity" name = "warehouse_quantity2" 
                                                    value = {inputs.warehouse_quantity2 || ""} onChange = {handleChange} ></input></td>
                                                </tr>
                                            </tbody>
                                        </table>
                                    <div className = "details-description">Product Variant Options</div>
                                        <table>
                                            <tbody>
                                                <tr>
                                                    <th style = {{textAlign: 'center'}}>Value</th>
                                                </tr>
                                                <tr>
                                                    <td><input type = '_text' style = {{width: '100%'}} placeholder = "Variant Option" name = "product_options" 
                                                    value = {inputs.product_options || ""} onChange = {handleChange} required></input></td>
                                                </tr>
                                                <tr>
                                                    <td><input type = '_text' style = {{width: '100%'}} placeholder = "Variant Option 2" name = "product_options2" 
                                                    value = {inputs.product_options2 || ""} onChange = {handleChange} required></input></td>
                                                </tr>
                                            </tbody>
                                        </table>
                                    
                                </div>
                            </div>
                            <div className = "details-right" style = {{backgroundColor: 'transparent'}}/>
                            <div className = "details-left" style = {{backgroundColor: 'transparent'}}/>
                        </div>
                    </div>
                </div>
                <button type = "submit" className = "submiit" style={{zIndex: '2'}}>Add Product</button>
                </form>
            </div>    
        </>
    );
};
  
Add_Product.defaultProps = 
{

};
export default Add_Product;