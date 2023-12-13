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
            document.getElementById(pageName).style.backgroundColor = "rgb(72, 101, 128)";
            document.getElementById(pageName).style.color = "black";
            
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
                <button className="tablink" id ="Variants">Addresses</button>
            </div>
        
            <div className="tabcontent" id="_Order" >
                <div className = "details-details">
                    <div className = "detailed" style={{backgroundColor: 'transparent'}}>
                        <div className = "details-title">Customer Orders</div>

                        <div className="order_header_div">
                            <div className="view_order_title">
                                {props.Order_Title}
                                <div className="view_order_status"></div>
                            </div>
                            <div>{props.Created_At}</div>
                        </div>
                        <div>
                            <table>
                                <tbody>
                                    <tr>
                                        <th>First Name</th>
                                        <th>Last Name</th>
                                        <th>Updated At</th>
                                    </tr>
                                    <tr>
                                        <td>{props.firstName}</td>
                                        <td>{props.lastName}</td>
                                        <td>{props.updatedAt}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>

                        <div className="order_lines_view">
                            <div className = "view_order_title">Order View</div>
                            <table className="order_table" style = {{left: '', transform: '', marginBottom: ''}}>
                                <tbody id = "detailed_table">

                                </tbody>
                            </table>
                        </div>
                        <div className="order_totals_view">
                        <div className = "view_order_title">Order Totals</div>
                            <table className="order_totals_table">
                                <tbody id = "detailed_table_view">
                                
                                </tbody>    
                            </table>
                        </div>
                          
                    </div>
                    <div className = "details-right" style={{backgroundColor: 'transparent'}}>
                        
                    </div>
                    <div id = "address" className = "details-left" style={{backgroundColor: 'transparent'}}></div>
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
    Order_Price: 'Price', 
    updatedAt: 'Updated'
}
export default Detailed_order;
