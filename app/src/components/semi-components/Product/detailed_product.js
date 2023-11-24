import {useEffect} from 'react';
import { useRef } from "react";
import Product_Variants from './product_variants';
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
            document.getElementById(pageName).style.backgroundColor = "#e7e7e7";
            
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

        let activity = document.querySelector(".details-title").innerHTML;
        let status = document.querySelector(".inactive");
        if(activity != "1")
        {
            status.className = "inactive";
        }
        else 
        {
            status.className = "activee";
        }
          
    }, []);

    return (
        
        <div id = "detailss" style = {{display: props.Display}}>
            <div className = "button-holder">
                <button className="tablink" id = "Product">Product</button>
                <button className="tablink" id ="Variants">Variants</button>
            </div>
        
            <div className="tabcontent" id="_Product" >
                <div className = "details-details">
                    <div className = "details-image" style = {{backgroundImage: `url(${props.Product_Image})`}}></div>
                    <div className = "detailed">
                        <div className = "details-title">{props.Product_Title}</div>
                        <i className = "inactive"/>
                        <span id = "activity">Activity</span>
                        <table>
                            <tbody>
                                <tr>
                                    <th>Product_Category</th>
                                    <th>Product_Code</th>
                                    <th>Product_Type</th>
                                    <th>Product_Price</th>
                                </tr>
                                <tr>
                                    <td>{props.Product_Category}</td>
                                    <td>{props.Product_Code}</td>
                                    <td>{props.Product_Type}</td>
                                    <td>{props.Product_Price}</td>
                                </tr>
                            </tbody>
                        </table> 
                        <div className = "details-description">Product Description</div>
                        <p>Product Description goes here, and it will be a extremelty long piece of text, of course, this can vary
                            but, on average, it could grow to be this large. But we'll see eyy!~</p>
                    </div>
                    <div className = "details-left"></div>
                    <div className = "details-right"></div>
                </div>
            </div>

            <div className="tabcontent" id="_Variants" >
                <div className = "details-details">
                <div className = "auto-slideshow-container" />
                    <div className = "detailed">
                        <div className = "details-title"> {props.Product_Title} Variants</div>
                        <div className = "variants" id="_variants" >
                        

                        </div>
                    </div>
                    <div className = "details-left"></div>
                    <div className = "details-right"></div>
                </div>
            </div>
        </div>
    );
};

Detailed_product.defaultProps = 
{
    Product_Title: 'Product title',
    Product_Code: 'Product code',
    Product_Options: 'Options',
    Product_Category: 'Category',
    Product_Type: 'Type',
    Product_Vendor: 'Vendor',
    Product_Image: '#ccc',
    Product_Price: 'Price',
    Variant_Title: 'Variant Title',
    Variant_Barcode: 'Variant Barcode',
    Variant_ID: 'Variant ID',
    Variant_SKU: 'Variant SKU',
    Variant_UpdateDate: 'Variant Update Date',
    Price_high: 'High Price',
    Price_low: 'Low Price'
}
export default Detailed_product;
/*
<div className = "details-image" style = {{backgroundImage: `linear-gradient(to bottom, transparent, white), url(${props.Product_Image})`}}>
*/