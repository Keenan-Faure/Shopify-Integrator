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
                pan[i].style.animation = "appear 0.4s ease-in";
            }
        }, 1400);

        /* Activity of pan elements */
        let activity = document.querySelectorAll(".p-d-activity");
        let option = document.querySelectorAll("#options");
        for(let i = 0; i < activity.length; i++)
        {
            if(activity[i].innerHTML == "1")
            {
                option[i].className = "p-d-true";
            }
            else if(activity[i].innerHTML == "")
            {
                option[i].className = "p-d-unknown";
            }
            else
            {
                option[i].className = "p-d-false";
            }
        }
        let category = document.querySelectorAll(".p-d-category"); let type = document.querySelectorAll(".p-d-type");
        let vendor = document.querySelectorAll(".p-d-vendor");

        for(let i = 0; i < category.length; i++)
        {
            if(category[i].innerHTML == "") { category[i].style.display = "none"; }
            if(type[i].innerHTML == "") { type[i].style.display = "none"; }
            if(vendor[i].innerHTML == "") { vendor[i].style.display = "none"; }
        }
        
    }, []);

    return (

        <div className = "pan">
            <div className = "pan-img" style = {{backgroundImage: `url(${props.Product_Image})`}}></div>
            <div className = "pan-details">
                <a className = "p-d-title">{props.Product_Title} <i id = "options" className = "p-d-options" title = "Activity"/></a> 
                <br/><br/>

                <a className = "p-d-code">{props.Product_Code}</a> <a className = "p-d-id">{props.Product_ID}</a> <a className = "p-d-activity">{props.Product_Activity}</a>
                <br/><br/>

                <a className = "p-d-category">{props.Product_Category}</a> <a className = "p-d-type">{props.Product_Type}</a> <a className = "p-d-vendor">{props.Product_Vendor}</a>
            </div>
        </div>
    );
};

Pan_details.defaultProps = 
{
    Product_Title: 'Product title',
    Product_Code: 'Product code',
    Product_Activity: '',
    Product_Category: 'Category',
    Product_Type: 'Type',
    Product_Vendor: 'Vendor',
}
export default Pan_details;