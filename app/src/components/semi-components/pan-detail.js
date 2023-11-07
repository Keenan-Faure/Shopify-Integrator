import {useEffect} from 'react';
import '../../CSS/page1.css';


function Pan_details(props)
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
        let pan_price = document.querySelectorAll(".pan-price");

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
                pan_price[i].style.color = "rgb(240, 248, 255, 0.8)";
            }
            pan[i].onmouseout = function(event)
            {
                let a_class = pan[i].querySelectorAll('a');
                for(let p = 0; p <a_class.length; p++)
                {
                    a_class[p].style.color = "black";
                }
                pan_details[i].style.color = "black";
                pan_price[i].style.color = "black"; 
            }
        }
    }, []);

    return (
        <div className = "pan">
            <div className = "pan-img"></div>
            <div className = "pan-details">
                <a href = "/#" className = "p-d-title">{props.Product_Title}</a>
                <br/><br/>

                <a href = "/#" className = "p-d-code">{props.Product_Code}</a>
                <br/><br/>

                <a href = "/#" className = "p-d-options">{props.Product_Options}</a> | <a href = "/#" className = "p-d-category">{props.Product_Category}</a> | <a href = "/#" className = "p-d-type">{props.Product_Type}</a> | <a href = "/#" className = "p-d-vendor">{props.Product_Vendor}</a>
            </div>
            <div className = "pan-price">
            {props.Product_Price}
            </div>
        </div>
    );
};

Pan_details.defaultProps = 
{
    Product_Title: 'Product title',
    Product_Code: 'Product code',
    Product_Options: 'Options',
    Product_Category: 'Category',
    Product_Type: 'Type',
    Product_Vendor: 'Vendor',
    Product_Price: 'Price'
}
export default Pan_details;