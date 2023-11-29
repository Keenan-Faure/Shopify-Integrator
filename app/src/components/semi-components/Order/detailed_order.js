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
                <button className="tablink" id = "Order">Order</button>
            </div>
        
            <div className="tabcontent" id="_Order" >
                <div className = "details-details">
                    <div className = "detailed-image" />
                    <div className = "detailed">
                        <div className = "details-title">Customer Orders</div>
                        <div className="order_lines_view">
                            <table className="order_table" style = {{left: '', transform: '', marginBottom: ''}}>
                                <tbody>
                                    <tr className="order_line">
                                        <td className="order_line_image"></td>
                                        <td className="order_line_title">
                                            <a>I am a title of a product</a>
                                            <p><b>SKU Code</b></p>
                                        </td>
                                        <td className="order_line_qty">x1</td>
                                        <td className="order_line_price">1500.00</td>
                                    </tr>
                                    <tr className="order_line">
                                        <td className="order_line_image"></td>
                                        <td className="order_line_title">
                                            <a>I am a title of a product</a>
                                            <p><b>SKU Code</b></p>
                                        </td>
                                        <td className="order_line_qty">x1</td>
                                        <td className="order_line_price">1500.00</td>
                                    </tr>
                                    <tr className="order_line">
                                        <td className="order_line_image"></td>
                                        <td className="order_line_title">
                                            <a>I am a title of a product</a>
                                            <p><b>SKU Code</b></p>
                                        </td>
                                        <td className="order_line_qty">x1</td>
                                        <td className="order_line_price">1500.00</td>
                                    </tr>
                                    <tr className="order_line">
                                        <td className="order_line_image"></td>
                                        <td className="order_line_title">
                                            <a>I am a title of a product</a>
                                            <p><b>SKU Code</b></p>
                                        </td>
                                        <td className="order_line_qty">x5</td>
                                        <td className="order_line_price">2,300,00.00</td>
                                    </tr>
                                    <tr className="order_line">
                                        <td className="order_line_image"></td>
                                        <td className="order_line_title">
                                            <a>I am a title of a product</a>
                                            <p><b>SKU Code</b></p>
                                        </td>
                                        <td className="order_line_qty">x5</td>
                                        <td className="order_line_price">2,300,00.00</td>
                                    </tr>
                                </tbody>
                            </table>
                        </div>
                        <div className="order_totals_view">
                            <table className="order_totals_table">
                                <tbody>
                                    <tr className="order_totals_line">
                                        <td className="order_totals_headers">
                                            Sub total
                                        </td>
                                        <td className="order_totals_middle"></td>
                                        <td className="order_totals_value">2,300,500</td>
                                    </tr>
                                    <tr className="order_totals_line">
                                        <td className="order_totals_headers">
                                            Tax
                                        </td>
                                        <td className="order_totals_middle">10%</td>
                                        <td className="order_totals_value">4,500</td>
                                    </tr>
                                    <tr className="order_totals_line">
                                        <td className="order_totals_headers">
                                            Shipping
                                        </td>
                                        <td className="order_totals_middle">Standard Shipping</td>
                                        <td className="order_totals_value">500.00</td>
                                    </tr>
                                    <tr className="order_totals_line">
                                        <td className="order_totals_headers">
                                            Total
                                        </td>
                                        <td className="order_totals_middle"></td>
                                        <td className="order_totals_value">2,350,500</td>
                                    </tr>
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
                            <div className="view_order_title_date">23 November 2023 +UTC 0000.000</div>
                        </div>
                        <div className="customer_data_container">
                            <div className="customer_data_row">
                                <div className="customer_data_header">First Name</div>
                                <textarea className="customer_data_tiles">Keenan</textarea>
                            </div>
                            <div className="customer_data_row">
                                <div className="customer_data_header">Last Name</div>
                                <textarea className="customer_data_tiles">Faure</textarea>
                            </div>
                            <div className="customer_data_row">
                                <div className="customer_data_header">Email</div>
                                <textarea className="customer_data_tiles">keenan@stock2shop.com</textarea>
                            </div>
                            <div className="customer_data_row">
                                <div className="customer_data_header">Phone number</div>
                                <textarea className="customer_data_tiles">0897665123</textarea>
                            </div>
                        </div>
                    </div>
                    <div className = "details-left">

                    <div className="order_data_customer_address_view">
                        <div className="customer_address_title">Shipping Address</div>
                        <hr />
                        <div className="customer_address_header">Address1</div>
                        <div className="customer_address_tiles">14 Tracy Close</div>
                        <div className="customer_address_header">Address2</div>
                        <div className="customer_address_tiles">Montrose park</div>
                        <div className="customer_address_header">City</div>
                        <div className="customer_address_tiles">Cape Town</div>
                        <div className="customer_address_header">Suburb</div>
                        <div className="customer_address_tiles">Mitchells Plain</div>
                        <div className="customer_address_header">Postal Code</div>
                        <div className="customer_address_tiles">7785</div>
                    </div>
                    <div className="order_data_customer_address_view">
                        <div className="customer_address_title">Billing Address</div>
                        <hr />
                        <div className="customer_address_header">Address1</div>
                        <div className="customer_address_tiles">14 Tracy Close</div>
                        <div className="customer_address_header">Address2</div>
                        <div className="customer_address_tiles">Montrose park</div>
                        <div className="customer_address_header">City</div>
                        <div className="customer_address_tiles">Cape Town</div>
                        <div className="customer_address_header">Suburb</div>
                        <div className="customer_address_tiles">Mitchells Plain</div>
                        <div className="customer_address_header">Postal Code</div>
                        <div className="customer_address_tiles">7785</div>
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
/*

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
*/

/*
<div class="background">
                <div class="order_data_customer_address_view">
                    <div class="customer_address_title">Shipping Address</div>
                    <hr />
                    <div class="customer_address_header">Address1</div>
                    <div class="customer_address_tiles">14 Tracy Close</div>
                    <div class="customer_address_header">Address2</div>
                    <div class="customer_address_tiles">Montrose park</div>
                    <div class="customer_address_header">City</div>
                    <div class="customer_address_tiles">Cape Town</div>
                    <div class="customer_address_header">Suburb</div>
                    <div class="customer_address_tiles">Mitchells Plain</div>
                    <div class="customer_address_header">Postal Code</div>
                    <div class="customer_address_tiles">7785</div>
                </div>
                <div class="order_data_customer_address_view" style= {{top: '50%'}}>
                    <div class="customer_address_title">Billing Address</div>
                    <hr />
                    <div class="customer_address_header">Address1</div>
                    <div class="customer_address_tiles">14 Tracy Close</div>
                    <div class="customer_address_header">Address2</div>
                    <div class="customer_address_tiles">Montrose park</div>
                    <div class="customer_address_header">City</div>
                    <div class="customer_address_tiles">Cape Town</div>
                    <div class="customer_address_header">Suburb</div>
                    <div class="customer_address_tiles">Mitchells Plain</div>
                    <div class="customer_address_header">Postal Code</div>
                    <div class="customer_address_tiles">7785</div>
                </div>

                <div class="view">
                    <div class="order_header_div">
                        <div class="view_order_title">
                            #13551 - Paid
                            <b style= {{fontSize: '26px'}}>•</b>
                            <div class="view_order_status"></div>
                        </div>
                        <div class="view_order_title_date">23 November 2023 +UTC 0000.000</div>
                    </div>
                    <div class="order_data_container">
                        <button class="customer_button_dd">
                            <div class="customer_button_right_arrow"></div>
                            Customer details
                        </button>
                    </div>
                    <div class="customer_data_container">
                        <div class="customer_data_row">
                            <div class="customer_data_header">First Name</div>
                            <textarea class="customer_data_tiles">Keenan</textarea>
                        </div>
                        <div class="customer_data_row">
                            <div class="customer_data_header">Last Name</div>
                            <textarea class="customer_data_tiles">Faure</textarea>
                        </div>
                        <div class="customer_data_row">
                            <div class="customer_data_header">Email</div>
                            <textarea class="customer_data_tiles">keenan@stock2shop.com</textarea>
                        </div>
                        <div class="customer_data_row">
                            <div class="customer_data_header">Phone number</div>
                            <textarea class="customer_data_tiles">0897665123</textarea>
                        </div>
                    </div>
                </div>
                <div class="order_lines_view">
                    <table class="order_table">
                        <tr class="order_line">
                            <td class="order_line_image"></td>
                            <td class="order_line_title">
                                <a>I am a title of a product</a>
                                <p><b>SKU Code</b></p>
                            </td>
                            <td class="order_line_qty">x1</td>
                            <td class="order_line_price">1500.00</td>
                        </tr>
                        <tr class="order_line">
                            <td class="order_line_image"></td>
                            <td class="order_line_title">
                                <a>I am a title of a product</a>
                                <p><b>SKU Code</b></p>
                            </td>
                            <td class="order_line_qty">x1</td>
                            <td class="order_line_price">1500.00</td>
                        </tr>
                        <tr class="order_line">
                            <td class="order_line_image"></td>
                            <td class="order_line_title">
                                <a>I am a title of a product</a>
                                <p><b>SKU Code</b></p>
                            </td>
                            <td class="order_line_qty">x1</td>
                            <td class="order_line_price">1500.00</td>
                        </tr>
                        <tr class="order_line">
                            <td class="order_line_image"></td>
                            <td class="order_line_title">
                                <a>I am a title of a product</a>
                                <p><b>SKU Code</b></p>
                            </td>
                            <td class="order_line_qty">x5</td>
                            <td class="order_line_price">2,300,00.00</td>
                        </tr>
                        <tr class="order_line">
                            <td class="order_line_image"></td>
                            <td class="order_line_title">
                                <a>I am a title of a product</a>
                                <p><b>SKU Code</b></p>
                            </td>
                            <td class="order_line_qty">x5</td>
                            <td class="order_line_price">2,300,00.00</td>
                        </tr>
                    </table>
                </div>
                <div class="order_totals_view">
                    <table class="order_totals_table">
                        <tr class="order_totals_line">
                            <td class="order_totals_headers">
                                Sub total
                            </td>
                            <td class="order_totals_middle"></td>
                            <td class="order_totals_value">2,300,500</td>
                        </tr>
                        <tr class="order_totals_line">
                            <td class="order_totals_headers">
                                Tax
                            </td>
                            <td class="order_totals_middle">10%</td>
                            <td class="order_totals_value">4,500</td>
                        </tr>
                        <tr class="order_totals_line">
                            <td class="order_totals_headers">
                                Shipping
                            </td>
                            <td class="order_totals_middle">Standard Shipping</td>
                            <td class="order_totals_value">500.00</td>
                        </tr>
                        <tr class="order_totals_line">
                            <td class="order_totals_headers">
                                Total
                            </td>
                            <td class="order_totals_middle"></td>
                            <td class="order_totals_value">2,350,500</td>
                        </tr>
                    </table>
                </div>
            </div>
*/