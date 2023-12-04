import {useEffect} from 'react';
import {useState} from "react";
import $ from 'jquery';
import '../../CSS/login.css';
import Background from '../../components/Background';

function Add_Customer()
{
    const[inputs, setInputs] = useState({});

    const handleChange = (event) =>
    {
        const name = event.target.name;
        const value = event.target.value;
        setInputs(values => ({...values, [name]: value}))
    }
    const AddCustomer = (event) =>
    {
        event.preventDefault();
        console.log(inputs);

        /*
        $.post("http://localhost:8080/api/login", JSON.stringify(inputs),[], 'json')
        .done(function( _data) 
        {
            console.log(_data);
        })
        .fail( function(xhr) 
        {
            alert(xhr.responseText);
        });
        */
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

        /* Rain Functions */
        var makeItRain = function() 
        {
            //clear out everything
            $('.rain').empty();
          
            var increment = 0;
            var drops = "";
            var backDrops = "";
          
            while (increment < 100) 
            {

                //couple random numbers to use for various randomizations
                //random number between 98 and 1
                var randoHundo = (Math.floor(Math.random() * (98 - 1 + 1) + 1));
                //random number between 5 and 2
                var randoFiver = (Math.floor(Math.random() * (5 - 2 + 1) + 2));
                //increment
                increment += randoFiver;
                //add in a new raindrop with various randomizations to certain CSS properties
                drops += '<div class="drop" style="left: ' + increment + '%; bottom: ' 
                + (randoFiver + randoFiver - 1 + 100) + '%; animation-delay: 0.' + randoHundo 
                + 's; animation-duration: 0.5' + randoHundo + 's;"><div class="stem" style="animation-delay: 0.' 
                + randoHundo + 's; animation-duration: 0.5' + randoHundo + 's;"></div><div class="splat" style="animation-delay: 0.' 
                + randoHundo + 's; animation-duration: 0.5' + randoHundo + 's;"></div></div>';
                
                backDrops += '<div class="drop" style="right: ' + increment + '%; bottom: ' 
                + (randoFiver + randoFiver - 1 + 100) + '%; animation-delay: 0.' + randoHundo 
                + 's; animation-duration: 0.5' + randoHundo + 's;"><div class="stem" style="animation-delay: 0.' 
                + randoHundo + 's; animation-duration: 0.5' + randoHundo + 's;"></div><div class="splat" style="animation-delay: 0.' 
                + randoHundo + 's; animation-duration: 0.5' + randoHundo + 's;"></div></div>';
            }
          
            $('.rain.front-row').append(drops);
            $('.rain.back-row').append(backDrops);
        }
          
        $('.splat-toggle.toggle').on('click', function() 
        {
            $('body').toggleClass('splat-toggle');
            $('.splat-toggle.toggle').toggleClass('active');
            makeItRain();
        });
          
        $('.back-row-toggle.toggle').on('click', function() 
        {
            $('body').toggleClass('back-row-toggle');
            $('.back-row-toggle.toggle').toggleClass('active');
            makeItRain();
        });
        
        $('.single-toggle.toggle').on('click', function() 
        {
            $('body').toggleClass('single-toggle');
            $('.single-toggle.toggle').toggleClass('active');
            makeItRain();
        });

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
        let navbar = document.getElementById("navbar");
        let details = document.getElementById("detailss");
        close.addEventListener("click", ()=> 
        {
            close.style.display = "none";
            details.style.animation = "Fadeout 0.5s ease-out";
            
            window.location.href = "/customers";
        });


    }, []);

    return (
        <>
            <Background />
            <div className = 'modal1' id = "model">
                <div className = "back-row-toggle splat-toggle">
                    <div className = "rain front-row"></div>
                    <div className = "rain back-row"></div>
                    <div className = "toggles">
                        <div className = "splat-toggle toggle active"></div>
                    </div>
                </div>

                <form className = 'modal-content' style ={{opacity: '1'}} method = 'post' onSubmit={(event) => AddCustomer(event)} autoComplete='off' id = 'form1' encType="multipart/form-data">
                <div id = "detailss">
                    <div className = 'rtn-button'></div>
                    <div className = "button-holder">
                        <button type = "button" className="tablink" id = "Customer">Customer</button>
                        <button type = "button" className="tablink" id ="Variants">Shipping</button>
                    </div>
                
                    <div className="tabcontent" id="_Customer" >
                        <div className = "details-details">
                            <div className = "detailed-image" />
                            <div className = "detailed">
                                <div className = "details-title">
                                    <input type = '_text' style ={{fontSize:'20px', width: '500px'}} placeholder = "Customer ID" name = "customer_id" value = {inputs.customer_id || ""}  
                                    onChange = {handleChange} required></input>
                                    </div>
                                <table>
                                    <tbody>
                                        <tr>
                                            <th>Customer Email</th>
                                            <th>Customer Firstname</th>
                                            <th>Customer Lastname</th>
                                            <th>Customer Phone</th>
                                        </tr>
                                        <tr>
                                            <td><input type = '_text' style = {{width: '150px'}} placeholder = "Customer Email" name = "customer_email" 
                                            value = {inputs.customer_email || ""} onChange = {handleChange} required></input></td>
                                            <td><input type = '_text' style = {{width: '150px'}} placeholder = "Customer First name" name = "customer_firstname" 
                                            value = {inputs.customer_firstname || ""} onChange = {handleChange} required></input></td>
                                            <td><input type = '_text' style = {{width: '150px'}} placeholder = "Customer Last name" name = "customer_lastname" 
                                            value = {inputs.customer_lastname || ""} onChange = {handleChange} required></input></td>
                                            <td><input type = '_text' style = {{width: '150px'}} placeholder = "Customer Phone" name = "customer_phone" 
                                            value = {inputs.customer_phone || ""} onChange = {handleChange} required></input></td>
                                        </tr>
                                    </tbody>
                                </table>  
                            </div>
                            <div className = "details-right"></div>
                            <div className = "details-left"></div>
                        </div>
                    </div>

                    <div className="tabcontent" id="_Variants" >
                        <div className = "details-details">
                        <div className = "detailed-image" />
                            <div className = "detailed">
                                <div className = "details-title">Shipping Details</div>
                                <div className = "variants" id="_variants"> 
                                <table>
                                    <tbody>
                                        <tr>
                                            <th>Shipping Address 1</th>
                                        </tr>
                                        <tr>
                                            <td><input type = '_text' style = {{width: '150px'}} placeholder =  "Shipping Address line 1" name = "shipping_address_1" 
                                            value = {inputs.shipping_address_1 || ""} onChange = {handleChange} required></input></td>
                                        </tr>

                                        <tr>
                                            <th>Shipping Address 2</th>
                                        </tr>
                                        <tr>
                                            <td><input type = '_text' style = {{width: '150px'}} placeholder =  "Shipping Address line 2" name = "shipping_address_2" 
                                            value = {inputs.shipping_address_2 || ""} onChange = {handleChange} required></input></td>
                                        </tr>

                                        <tr>
                                            <th>Shipping Address 3</th>
                                        </tr>
                                        <tr>
                                            <td><input type = '_text' style = {{width: '150px'}} placeholder =  "Shipping Address line 3" name = "shipping_address_3" 
                                            value = {inputs.shipping_address_3 || ""} onChange = {handleChange} required></input></td>
                                        </tr>

                                        <tr>
                                            <th>Shipping Address 4</th>
                                        </tr>
                                        <tr>
                                            <td><input type = '_text' style = {{width: '150px'}} placeholder =  "Shipping Address line 4" name = "shipping_address_4" 
                                            value = {inputs.shipping_address_4 || ""} onChange = {handleChange} required></input></td>
                                        </tr>
                                        <tr>
                                            <th>Shipping Address 5</th>
                                        </tr>
                                        <tr>
                                            <td><input type = '_text' style = {{width: '150px'}} placeholder =  "Shipping Address line 5" name = "shipping_address_5" 
                                            value = {inputs.shipping_address_5 || ""} onChange = {handleChange} required></input></td>
                                        </tr>
                                    </tbody>
                                </table> 
                                </div>
                            </div>
                            <div className = "details-right"></div>
                            <div className = "details-left"></div>
                        </div>
                    </div>
                </div>
                <button type = "submit" className = "submiit">Add Customer</button>           
                </form>
            </div>    
        </>
    );
};
  
Add_Customer.defaultProps = 
{

};
export default Add_Customer;