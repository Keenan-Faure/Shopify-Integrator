import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_customer(props)
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

        let home = document.getElementById("Customer");
        home.addEventListener("click", () =>
        {
            openPage('Customer');
        });

        let defaul = document.getElementById("Variants");
        defaul.addEventListener("click", () =>
        {
            openPage('Variants');
        });

        document.getElementById("Customer").click();

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
                <button className="tablink" id = "Customer">Customer</button>
                <button className="tablink" id ="Variants">Variants</button>
            </div>
        
            <div className="tabcontent" id="_Customer" >
                <div className = "details-details">
                    <div className = "detailed-image" />
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
                        <div className = "details-description">Customer Descriptions</div>
                        <p id = "description">
                            Customer Description goes here, and it will be a extremelty long piece of text, of course, this can vary
                            but, on average, it could grow to be this large. But we'll see eyy!~
                        </p>
                        <div className = "details-description">Customer Warehousing</div> 
                        <div className = "details-warehousing"></div>  
                    </div>
                    <div className = "details-left"></div>
                    <div className = "details-right"></div>
                </div>
            </div>

            <div className="tabcontent" id="_Variants" >
                <div className = "details-details">
                <div className = "detailed-image" />
                    <div className = "detailed">
                        <div className = "details-title"> {props.Product_Title} Variants</div>
                        <div className = "variants" id="_variants" ></div>
                    </div>
                    <div className = "details-left"></div>
                    <div className = "details-right"></div>
                </div>
            </div>
        </div>
    );
};

Detailed_customer.defaultProps = 
{
    Customer_Title: 'Customer title',
    Customer_Code: 'Customer code',
    Customer_Options: 'Options',
    Customert_Category: 'Category',
    Customer_Type: 'Type',
    Customer_Vendor: 'Vendor',
    Customer_Price: 'Price'
}
export default Detailed_customer;