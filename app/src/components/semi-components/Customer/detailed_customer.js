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
            document.getElementById(pageName).style.backgroundColor = "rgb(72, 101, 128)";
            document.getElementById(pageName).style.color = "black";
            
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
                navbar.style.animation = "FadeIn ease-in 0.6s";
                details.style.display = "none";
                navbar.style.display = "block";
                main.style.display = "block";
            }, 500);
        });

    }, []);

    return (
        
        <div id = "detailss" style = {{display: props.Display}}>
            <div className = 'rtn-button'></div>
            <div className = "button-holder">
                <button className="tablink" id = "Customer">Customer</button>
                <button className="tablink" id ="Variants">Shipping</button>
            </div>
        
            <div className="tabcontent" id="_Customer" >
                <div className = "details-details">
                    <div className = "detailed-image" />
                    <div className = "detailed" style={{backgroundColor: 'transparent'}}>
                        <div className = "details-title">{props.Customer_ID}</div>
                        <table>
                            <tbody>
                                <tr>
                                    <th>Customer Email</th>
                                    <th>Customer Firstname</th>
                                    <th>Customer Lastname</th>
                                    <th>Customer Phone</th>
                                </tr>
                                <tr>
                                    <td>{props.Customer_Email}</td>
                                    <td>{props.Customer_firstName}</td>
                                    <td>{props.Customer_lastName}</td>
                                    <td>{props.Customer_Phone}</td>
                                </tr>
                            </tbody>
                        </table> 
                        <div className = "details-description">Customer Warehousing</div> 
                        <div className = "details-warehousing"></div>  
                    </div>
                    <div className = "details-right" style={{backgroundColor: 'transparent'}}></div>
                    <div className = "details-left" style={{backgroundColor: 'transparent'}}></div>
                </div>
            </div>

            <div className="tabcontent" id="_Variants" >
                <div className = "details-details">
                <div className = "detailed-image" />
                    <div className = "detailed" style={{backgroundColor: 'transparent'}}>
                        <div className = "details-title"> {props.Customer_firstName} Shipping Addresses</div>
                        <div className = "variants" id="_variants" ></div>
                    </div>
                    <div className = "details-right" style={{backgroundColor: 'transparent'}}></div>
                    <div className = "details-left" style={{backgroundColor: 'transparent'}}></div>
                </div>
            </div>
        </div>
    );
};

Detailed_customer.defaultProps = 
{
    Customer_ID: 'Customer ID',
    Customer_Email: 'Customer code',
    Customer_firstName: 'Options',
    Customer_lastName: 'Category',
    Customer_Phone: 'Type',
    Customer_Vendor: 'Vendor',
    Customer_Updated: 'Price'
}
export default Detailed_customer;