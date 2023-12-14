import {useEffect} from 'react';
import '../../CSS/page1.css';


function Queue_details(props)
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

        /* Hover brightens the color of the pan element details */
        let pan_details = document.querySelectorAll(".pan-details");
        for(let i = 0; i < pan.length ; i++)
        {
            pan[i].onmouseover = function(event)
            {
                let a_class = pan[i].querySelectorAll('a');
                for(let p = 0; p < a_class.length; p++)
                {
                    a_class[p].style.color = "rgb(240, 248, 255, 0.8)"
                }
                pan_details[i].style.color = "rgb(240, 248, 255, 0.8)";
            }
            pan[i].onmouseout = function(event)
            {
                let a_class = pan[i].querySelectorAll('a');
                for(let p = 0; p < a_class.length - 1; p++)
                {
                    a_class[p].style.color = "black";
                }
                pan_details[i].style.color = "black";


                for(let p = a_class.length - 1; p < a_class.length; p++)
                {
                    a_class[p].style.color = "#afafaf";
                }
                pan_details[i].style.color = "#afafaf";
            }
        }


        /* Activity of pan elements */
        let activity = document.querySelectorAll(".status");
        let option = document.querySelectorAll("#options");
        for(let i = 0; i < activity.length; i++)
        {
            if(activity[i].innerHTML == "in-queue")
            {
                option[i].className = "p-d-load";
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
    }, []);

    return (
        <div className = "pan" style = {{height: '90px'}}>
            <div className = "pan-img" style={{height: '90px', width: '115px'}}></div>
            <div className = "pan-details">
                <a href = "/#" className = "p-d-title"><b>Queue Title:  </b>{props.Queue_Instruction} <i id = "options" href = "/#" className = "p-d-options" /></a> 
                <br/><br/>

                <a href = "/#" className = "p-d-type"><b>Queue Type:  </b>{props.Queue_Type}</a> <a href = "/#" className = "p-d-id">{props.Queue_ID}</a> <a href = "/#" className = "p-d-status" style={{color: '#afafaf'}}><b>Queue Status:</b>  <div className = "status">{props.Queue_Status}</div></a> 
            </div>
        </div>
    );
};

Queue_details.defaultProps = 
{
    Queue_Creation_Date: 'N/A',
    Queue_Instruction: 'N/A',
    Queue_Updated_At: 'N/A',
    Queue_Type: 'N/A',
    Queue_Status: 'N/A',
    Queue_ID: 'N/A',
}
export default Queue_details;