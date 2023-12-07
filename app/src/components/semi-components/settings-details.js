import {useEffect} from 'react';
import '../../CSS/page2.css';

function Setting_details(props)
{
    useEffect(()=> 
    {
        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        window.onload = function(event)
        {
            navigation.style.left = "0%";
            navigation.style.position = "relative";
            navigation.style.width = "100%";
        }
        
    }, []);

    return (
        <>
            <div className = "setting">
                <div className = "setting-title">{props.subTitle1}</div>
                <div className = "setting-details description">{props.Description1}</div>
                <div className = "setting-details value">{props.Value1}</div>
                <button className = "button-on-off">Turn on</button>
            </div>
        </>
    );
}

Setting_details.defaultProps = 
{
    subTitle1: 'Sub Title of setting',
    Description1: 'Description of product goes here, as well as any additional information',
    Value1: 'Value of setting currently in the api',
}

export default Setting_details;