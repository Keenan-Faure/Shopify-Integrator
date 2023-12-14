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
        <div className = "pan" >
            <div className = "pan-img img-3"></div>
            <div className = "pan-details">
                <a href = "/#" className = "p-d-title">{props.Customer_ID}</a>
                <br/><br/>

                <a href = "/#" className = "p-d-code">{props.Customer_firstName}</a>
                <br/><br/>

                <a href = "/#" className = "p-d-options">{props.Customer_lastName}</a>
            </div>
        </div>
    );
};

Customer_details.defaultProps = 
{
    Customer_ID: 'Customer ID',
    Customer_firstName: 'Customer firstName',
    Customer_lastName: 'Customer lastName',
}
export default Customer_details;