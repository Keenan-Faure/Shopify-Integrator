import {useEffect} from 'react';
import '../../CSS/page1.css';


function Customer_details(props)
{
    useEffect(()=> 
    {
        let pan = document.querySelectorAll(".pan");
        setTimeout(() =>
        {
            for(let i = 0; i < pan.length; i ++)
            {
                pan[i].style.display = "block";
                pan[i].style.animation = "appear 1.2s ease-in";
            }
        }, 1400);

        /* Hover brightens the color of the pan element details */
        let pan_details = document.querySelectorAll(".pan-details");

        for(let i = 0; i < pan.length; i++)
        {
            pan[i].onmouseover = function(event)
            {
                let a_class = pan[i].querySelectorAll('a');
                for(let p = 0; p <a_class.length; p++)
                {
                    a_class[p].style.color = "rgb(240, 248, 255, 0.8)"
                }
                pan_details[i].style.color = "rgb(240, 248, 255, 0.8)";
            }
            pan[i].onmouseout = function(event)
            {
                let a_class = pan[i].querySelectorAll('a');
                for(let p = 0; p <a_class.length; p++)
                {
                    a_class[p].style.color = "black";
                }
                pan_details[i].style.color = "black";
            }
        }
    }, []);

    return (
        <div className = "pan">
            <div className = "pan-img"></div>
            <div className = "pan-details">
                <a href = "/#" className = "p-d-title">{props.Customer_Title}</a>
                <br/><br/>

                <a href = "/#" className = "p-d-code">{props.Customer_Code}</a>
                <br/><br/>

                <a href = "/#" className = "p-d-options">{props.Customer_Options}</a> | <a href = "/#" className = "p-d-category">{props.Customer_Category}</a> | <a href = "/#" className = "p-d-type">{props.Customer_Type}</a> | <a href = "/#" className = "p-d-vendor">{props.Customer_Vendor}</a>
            </div>
        </div>
    );
};

Customer_details.defaultProps = 
{
    Customer_Title: 'Customer title',
    Customer_Code: 'Customer code',
    Customer_Options: 'Options',
    Customer_Category: 'Category',
    Customer_Type: 'Type',
    Customer_Vendor: 'Vendor',
    Customer_Price: 'Price'
}
export default Customer_details;