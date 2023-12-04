import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_order(props)
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

        let home = document.getElementById("Order");
        home.addEventListener("click", () =>
        {
            openPage('Order');
        });

        document.getElementById("Order").click();

        /* When the user clicks on the return button */
        let close = document.querySelector(".rtn-button");
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
                <button className="tablink" id = "Order">Order</button>
            </div>
        
            <div className="tabcontent" id="_Order" >
                <div className = "details-details">
                    <div className = "detailed-image" />
                    <div className = "detailed">
                        <div className = "details-title">Customer Orders</div>
                        <div className="order_lines_view">
                            <table className="order_table" style = {{left: '', transform: '', marginBottom: ''}}>
                                <tbody id = "detailed_table">
                                    
                                </tbody>
                            </table>
                        </div>
                        <div className="order_totals_view">
                            <table className="order_totals_table">
                                <tbody id = "detailed_table_view">
                                
                                </tbody>    
                            </table>
                        </div>
                          
                    </div>
                    <div className = "details-right">
                        <div className="order_header_div">
                            <div className="view_order_title">
                                #13551 - Paid
                                <b style= {{fontSize: '26px'}}>•</b>
                                <div className="view_order_status"></div>
                            </div>
                            <div className="view_order_title_date">{props.Created_At}</div>
                        </div>
                        <div className="customer_data_container">
                            <div className="customer_data_row">
                                <div className="customer_data_header">First Name</div>
                                <div className="customer_data_tiles">{props.firstName}</div>
                            </div>
                            <div className="customer_data_row">
                                <div className="customer_data_header">Last Name</div>
                                <div className="customer_data_tiles">{props.lastName}</div>
                            </div>
                            <div className="customer_data_row">
                                <div className="customer_data_header">Email</div>
                                <div className="customer_data_tiles">{props.shippingAddress}</div>
                            </div>
                            <div className="customer_data_row">
                                <div className="customer_data_header">Phone number</div>
                                <div className="customer_data_tiles">{props.PhoneNum}</div>
                            </div>
                        </div>
                    </div>
                    <div className = "details-left">

                    <div className="order_data_customer_address_view">
                        <div className="customer_address_title">Shipping Address</div>
                        <hr />
                        <div className="customer_address_header">Address1</div>
                        <div className="customer_address_tiles">{props.S_Address1}</div>
                        <div className="customer_address_header">Address2</div>
                        <div className="customer_address_tiles">{props.S_Address2}</div>
                        <div className="customer_address_header">City</div>
                        <div className="customer_address_tiles">{props.S_Address3}</div>
                        <div className="customer_address_header">Suburb</div>
                        <div className="customer_address_tiles">{props.S_Address4}</div>
                        <div className="customer_address_header">Postal Code</div>
                        <div className="customer_address_tiles">{props.S_Address5}</div>
                    </div>
                    <div className="order_data_customer_address_view">
                        <div className="customer_address_title">Billing Address</div>
                        <hr />
                        <div className="customer_address_header">Address1</div>
                        <div className="customer_address_tiles">{props.B_Address1}</div>
                        <div className="customer_address_header">Address2</div>
                        <div className="customer_address_tiles">{props.B_Address2}</div>
                        <div className="customer_address_header">City</div>
                        <div className="customer_address_tiles">{props.B_Address3}</div>
                        <div className="customer_address_header">Suburb</div>
                        <div className="customer_address_tiles">{props.B_Address4}</div>
                        <div className="customer_address_header">Postal Code</div>
                        <div className="customer_address_tiles">{props.B_Address5}</div>
                    </div>

                    </div>
                </div>
            </div>
        </div>
    );
};

Detailed_order.defaultProps = 
{
    Order_Title: 'Order title',
    Order_Code: 'Order code',
    Order_Options: 'Options',
    Order_Category: 'Category',
    Order_Type: 'Type',
    Order_Vendor: 'Vendor',
    Order_Price: 'Price'
}
export default Detailed_order;
