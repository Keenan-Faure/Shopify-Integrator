import {useEffect} from 'react';
import '../../CSS/page1.css';

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
            document.getElementById(pageName).style.backgroundColor = "#b6b6b6";
            
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

        let contact = document.getElementById("Images");
        contact.addEventListener("click", () =>
        {
            openPage('Images');
        });

        document.getElementById("Product").click();
          
    }, []);

    return (
        <div>
            <div className = "button-holder">
                <button className="tablink" id = "Product">Product</button>
                <button className="tablink" id ="Variants">Variants</button>
                <button className="tablink" id = "Images">Images</button>
            </div>
        
            <div id="_Product" className="tabcontent">
                <div className = "details-details">
                    <div className = "details-image" style = {{backgroundImage: `url(${props.Product_Image})`}}></div>
                    <div className = "detailed">
                        <div className = "details-title">
                            <i className = "active"/>{props.Product_Title}
                            <span id = "activity">Activity</span>
                        </div>
                        <table>
                            <tbody>
                                <tr>
                                    <th>Product_Category</th>
                                    <th>Product_Code</th>
                                    <th>Product_Type</th>
                                    <th>Product_Price</th>
                                </tr>
                                <tr>
                                    <td><div className = "details-category">Product_Category</div></td>
                                    <td><div className = "details-code">Product_Code</div></td>
                                    <td><div className = "details-type">Product_Type</div></td>
                                    <td><div className = "details-price">Product_Price</div></td>
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

            <div id="_Variants" className="tabcontent">
            <div className = "details-details">
                    <div className = "details-image" style = {{backgroundImage: `url(${props.Product_Image})`}}></div>
                    <div className = "detailed">
                        
                    </div>
                    <div className = "details-left"></div>
                    <div className = "details-right"></div>
                </div>
            </div>

            <div id="_Images" className="tabcontent">
            <div className = "details-details">
                    <div className = "details-image" style = {{backgroundImage: `url(${props.Product_Image})`}}></div>
                    <div className = "detailed">
                        
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
    Product_Price: 'Price'
}
export default Detailed_product;
/*
<div className = "details-image" style = {{backgroundImage: `linear-gradient(to bottom, transparent, white), url(${props.Product_Image})`}}>
*/