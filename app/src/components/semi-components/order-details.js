import {useEffect} from 'react';
import '../../CSS/page1.css';


function Order_details(props)
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
                <a href = "/#" className = "p-d-title">{props.Order_Title}</a>
                <br/><br/>

                <a href = "/#" className = "p-d-code">{props.Order_Code}</a>
                <br/><br/>

                <a href = "/#" className = "p-d-options">{props.Order_Options}</a> | <a href = "/#" className = "p-d-category">{props.Order_Category}</a> | <a href = "/#" className = "p-d-type">{props.Order_Type}</a> | <a href = "/#" className = "p-d-vendor">{props.Order_Vendor}</a>
            </div>
        </div>
    );
};

Order_details.defaultProps = 
{
    Order_Title: 'Order title',
    Order_Code: 'Order code',
    Order_Options: 'Options',
    Order_Category: 'Category',
    Order_Type: 'Type',
    Order_Vendor: 'Vendor',
    Order_Price: 'Price'
}
export default Order_details;