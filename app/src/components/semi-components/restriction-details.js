import {useEffect} from 'react';
import '../../CSS/page2.css';

function Restriction_details(props)
{
    useEffect(()=> 
    {
        /* Activity of  elements */

       let activit = document.querySelectorAll(".flag_value");
       let slider_ = document.querySelectorAll("#lide");
       for(let i = 0; i <activit.length; i++)
       {
            if(activit[i].innerHTML == "app")
            {
                slider_[i].checked = true;
            }
            else if(activit[i].innerHTML == "")
            {
                slider_[i].checked = false;
            }
            else
            {
                slider_[i].checked = false;
            }
       }
       

    }, []);

    return (
        <>
            <table className = "table" id = "_table" style = {{width: '40%', marginTop: '10px'}}>
                <tbody>
                    <tr>
                        <th style = {{width: '33%'}}>Restriction Name</th>
                        <th style = {{width: '33%'}}>Flag</th>
                        <th style = {{width: '33%'}}>Current Value</th>
                    </tr>
                    <>
                        {props.restriction}
                    </>
                </tbody>
            </table>
            
        </>
    );
}

Restriction_details.defaultProps = 
{
    Key: 'Sub Title of setting',
    Description: 'Description of product goes here, as well as any additional information',
    Value: 'Value of setting currently in the api',
    id: 'id of the setting'
}

export default Restriction_details;