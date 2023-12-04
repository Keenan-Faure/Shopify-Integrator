import {useEffect} from 'react';
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

        /* When the user clicks on the return button */
        let close = document.querySelector(".rtn-button");
        let filter = document.querySelector(".filter");
        let main = document.querySelector(".main");
        let navbar = document.getElementById("navbar");
        let details = document.querySelector(".details");
        close.addEventListener("click", ()=> 
        {
            close.style.display = "none";
            details.style.animation = "Fadeout 0.5s ease-out";
            
            setTimeout(() => 
            {
                main.style.animation = "FadeIn ease-in 0.6s";
                filter.style.animation = "FadeIn ease-in 0.6s";
                navbar.style.animation = "FadeIn ease-in 0.6s";
                details.style.display = "none";
                navbar.style.display = "block";
                main.style.display = "block";
                filter.style.display = "block";
            }, 500);
        });

    }, []);

    return (
        
        <div id = "detailss" style = {{display: props.Display}}>
            <div className = 'rtn-button'></div>
            <div className = "button-holder">
                <button className="tablink" id = "Product">Product</button>
                <button className="tablink" id ="Variants">Variants</button>
            </div>
        
            <div className="tabcontent" id="_Product" >
                <div className = "details-details">
                    <div className = "auto-slideshow-container" />
                    <div className = "detailed">
                        <div className = "details-title">{props.Product_Title}<i className = "inactive"/></div>
                        
                        <span id = "activity">Activity</span>
                        <table>
                            <tbody>
                                <tr>
                                    <th>Product_Category</th>
                                    <th>Product_Code</th>
                                    <th>Product_Type</th>
                                    <th>Product_Vendor</th>
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
                        <p id = "description">
                            Product Description goes here, and it will be a extremelty long piece of text, of course, this can vary
                            but, on average, it could grow to be this large. But we'll see eyy!~
                        </p>
                        <div className = "details-description">Product Warehousing</div> 
                        <div className = "details-warehousing"></div>  
                    </div>
                    <div className = "details-right"></div>
                    <div className = "details-left"></div>
                </div>
            </div>

            <div className="tabcontent" id="_Variants" >
                <div className = "details-details">
                    <div className = "auto-slideshow-container" />
                    <div className = "detailed">
                        <div className = "details-title"> {props.Product_Title} Variants</div>
                        <div className = "variants" id="_variants" ></div>
                    </div>
                    <div className = "details-right"></div>
                    <div className = "details-left"></div>
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
}
export default Detailed_product;