import {useEffect} from 'react';
import $ from 'jquery';
import '../../../CSS/detailed.css';

function Detailed_product(props)
{
    useEffect(()=> 
    {
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
        home.addEventListener("click", () => { openPage('Product'); });

        let defaul = document.getElementById("Variants");
        defaul.addEventListener("click", () => { openPage('Variants'); });
        document.getElementById("Product").click();

        /* Activity of the product */
        let activity = document.querySelector(".details-title").innerHTML;
        let status = document.querySelector(".inactive");
        if(activity != "1") { status.className = "inactive"; }
        else { status.className = "activee"; }

        /* When the user clicks on the return button */
        let close = document.querySelector(".rtn-button");
        let filter = document.querySelector(".filter");
        let main = document.querySelector(".main");
        let navbar = document.getElementById("navbar");
        let details = document.querySelector(".details");
        close.addEventListener("click", ()=> 
        {
            close.style.display = "none";
            details.style.animation = "Fadeout 1s ease-out";
            setTimeout(() => 
            {
                main.style.animation = "FadeIn ease-in 1s";
                filter.style.animation = "FadeIn ease-in 1s";
                navbar.style.animation = "FadeIn ease-in 1s";
                details.style.display = "none";
                navbar.style.display = "block";
                main.style.display = "block";
                filter.style.display = "block";
            }, 500);
        });

        /* Edit Feature */
        let edit = document.getElementById("Edit");
        let confirm = document.querySelector(".confirm-line");
        edit.addEventListener("click", () =>
        {
            let td_list = document.querySelectorAll("td"); let description = document.getElementById("description");
            let variant_updateDate = document.querySelector(".variant-updateDate");
            confirm.style.display = "block";
            for(let i = 0; i< td_list.length; i++)
            {
                td_list[i].contentEditable = "true";
            }
            description.contentEditable = "true"; variant_updateDate.contentEditable = "true";
            
        });

        confirm.addEventListener("click", () =>
        {
            let td_list = document.querySelectorAll("td"); let description = document.getElementById("description");
            let variant_updateDate = document.querySelector(".variant-updateDate"); let price = document.querySelectorAll(".price");
            let barcode = document.querySelectorAll(".barcode"); let sku = document.querySelectorAll(".sku"); 
            let option1 = document.querySelectorAll(".option1"); let option2 = document.querySelectorAll(".option2"); let option3 = document.querySelectorAll(".option3");
            confirm.style.display = "none";

            let title = document.getElementById("title");
            for(let i = 0; i< td_list.length; i++)
            {
                td_list[i].contentEditable = "false";
            }
            description.contentEditable = "false"; variant_updateDate.contentEditable = "false";

            let object = 
            {
                body_html: description.innerHTML, 
                category: td_list[0].innerHTML, 
                options:
                [

                ],
                product_code: td_list[1].innerHTML, 
                product_type: td_list[2].innerHTML, 
                title: title.innerHTML, 
                variants: 
                [

                ],
                vendor: td_list[3].innerHTML
            };
            let quantities = {};
            let options = {};
            
            let _quantities = document.querySelectorAll(".quantities");
            let price_name = document.querySelectorAll(".price_name");
            let price_value = document.querySelectorAll(".price_value");
            let _options = document.querySelectorAll(".product_options");

            /* Options object */
            for(let i = 0; i < _options.length; i++)
            {
                if(_options[i].childNodes.length <= 1)
                {
                    options =
                    {         
                        name: "",
                        value: ""
                    };
                }
                else 
                {
                    options =
                    {
                        name: _options[i].childNodes[0].innerHTML,
                        value: _options[i].childNodes[1].innerHTML
                    };
                }
                object.options[i] = options;
            }

            /* Quantities object + rest of it */
            for(let i = 0; i < price.length; i++)
            {
                /* Keep variants variable inside, so it can start fresh when the for loop restarts */
                let variants = {};

                if(_quantities[i].childNodes.length <= 1)
                {
                    quantities =
                    [{
                        
                        name: "",
                        value: 0
                    }];
                }
                else 
                {
                    quantities =
                    [{
                        name: _quantities[i].childNodes[0].innerHTML,
                        value: parseInt(_quantities[i].childNodes[1].innerHTML)
                    }];
                }
                variants.sku = sku[i].innerHTML; 
                variants.barcode = barcode[i].innerHTML;
                variants.option1 = option1[i].innerHTML; 
                variants.option2 = option2[i].innerHTML; 
                variants.option3 = option3[i].innerHTML;

                variants.variant_quantities = quantities;
                variants.variant_price_tiers = 
                [{

                
                    name: price_name[i].innerHTML,
                    value: price_value[i].innerHTML
                }];
                object.variants[i] = variants;
            }
            
            let id = document.querySelector("._id").innerHTML;
            console.log(object);
            
            const api_key = localStorage.getItem('api_key');
            $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key}, type: 'PUT' });

            $.ajax({ type: 'PUT', url: "http://localhost:8080/api/products/" + id, 
            contentType: 'application/json', data: JSON.stringify(object)})
            .done(function (_data) 
            {
                console.log(_data);
            })
            .fail( function(xhr) 
            {
                alert(xhr.responseText);
            });

        })

    }, []);

    return (
        
        <div id = "detailss" style = {{display: props.Display}}>
            <div className = 'rtn-button'></div>
            <div className = "button-holder" style = {{position: 'absolute', width: '71%', left:'29%'}}>
                <button className="tablink" id = "Product" style ={{left: '-14%', width:'95px'}}>Product</button>
                <button className="tablink" id ="Variants" style ={{left: '-14%', width:'95px'}}>Variants</button>
                <button className="tablink" id = "Edit" style ={{left: '-14%', width:'95px'}}>Edit</button>
            </div>
        
            <div className="tabcontent" id="_Product" >
                <div className = "details-details">
                    <div className = "auto-slideshow-container" style={{backgroundColor: 'transparent'}} />
                    <div className = "detailed">
                        <div className = "details-title"><div id ="_title" style={{position: 'relative',top: '10px',display: 'inline-block'}}>
                            <div id = "title">{props.Product_Title}</div></div>
                            <i className = "inactive"/>
                        </div>
                        <div className = "_id" style ={{display: 'none'}}>{props.Product_ID}</div>
                        
                        <span id = "activity">Activity</span>
                        <table>
                            <tbody>
                                <tr>
                                    <th>Product Category</th>
                                    <th>Product Code</th>
                                    <th>Product Type</th>
                                    <th>Product Vendor</th>
                                </tr>
                                <tr>
                                    <td>{props.Product_Category}</td>
                                    <td>{props.Product_Code}</td>
                                    <td>{props.Product_Type}</td>
                                    <td>{props.Product_Vendor}</td>
                                </tr>
                            </tbody>
                        </table> 
                        <div className = "details-description">Product Description</div>
                        <div className = "description" id = "description" style = {{resize:'none'}} rows = "5" cols = "80">{props.Product_Description}</div>

                        <div className = "details-description">Product Options</div> 
                        <div className = "details-options" style={{marginBottom: '15px'}}>
                            <table>
                                <tbody>
                                    <tr>
                                        <th style= {{width: '50%'}}>Value</th>
                                        <th style= {{width: '50%'}}>Position</th>
                                    </tr>
                                    <>
                                        {props.Product_Options}
                                    </>
                                </tbody>
                            </table>
                        </div>
                    </div>
                    <div className = "details-right" style ={{backgroundColor: 'transparent'}}></div>
                    <div className = "details-left"></div>
                </div>
            </div>

            <div className="tabcontent" id="_Variants" >
                <div className = "details-details">
                    <div className = "auto-slideshow-container" style={{backgroundColor: 'transparent'}} />
                    <div className = "detailed">
                        <div className = "details-title"> {props.Product_Title} Variants</div>
                        <div className = "variants" id="_variants" ></div>
                    </div>
                    <div className = "details-right" style ={{backgroundColor: 'transparent'}}></div>
                    <div className = "details-left"></div>
                </div>
            </div>
            <div className = "confirm-line">
                <button className="tablink" id = "confirm" style ={{left: '50%'/*, transform: 'translate(-50%)'*/}}>Save</button>
            </div>
        </div>
    );
};

Detailed_product.defaultProps = 
{
    Product_Title: 'N/A',
    Product_Code: 'N/A',
    Product_Options: 'N/A',
    Product_Category: 'N/A',
    Product_Type: 'N/A',
    Product_Vendor: 'N/A',
    Product_Image: 'N/A',
    Product_Price: 'N/A',
    Product_Options: 'N/A'
}
export default Detailed_product;