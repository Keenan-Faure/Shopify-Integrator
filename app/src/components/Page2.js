import {useEffect} from 'react';
import Background from './Background';

import '../CSS/page2.css';


function Page2(props)
{
    useEffect(()=> 
    {
        /* Ensure the model is shown */
        let model = document.getElementById("model");
        let navbar = document.getElementById("navbar");
        navbar.style.display = "block";
        model.style.display = "none";

        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        window.onload = function(event)
        {
            navigation.style.left = "30%";
            navigation.style.position = "absolute";
            navigation.style.width = "70%";
            navigation.style.animation = "MoveLeft 1.2s ease";
        }
    }, []);

    return (
        <>
            <Background />
            <div className = "component1">
                <div className = "main-container">
                    <div className = "settings">
                        <div className = "app-settings">
                            <div className = "title">{props.Title1}</div>
                            <div className = "setting">
                                <div className = "setting-title">{props.subTitle1}</div>
                                <div className = "setting-details description">{props.Description1}</div>
                                <div className = "setting-details value">{props.Value1}</div>
                                <div className = "button-on-off">Turn on </div>
                            </div>
                        </div>
                        <div className = "shopify-settings">
                            <div className = "title">{props.Title2}</div>
                            <div className = "setting">
                                <div className = "setting-title">{props.subTitle1}</div>
                                <div className = "setting-details description">{props.Description1}</div>
                                <div className = "setting-details value">{props.Value1}</div>
                                <div className = "button-on-off">Turn on </div>
                            </div>
                        </div>
                    </div>    
                </div>
                <div className = "side-container">
                    <div className = "settings-2">
                        <div className = "application"><i className = "a"/>Application Settings:</div>
                        <div className = "mini-setting">Setting1</div>
                        <div className = "mini-setting">Setting2</div>

                    </div>
                    <div className = "settings-2">
                        <div className = "application"><i className = "b"/>Spotify Settings:</div>
                        <div className = "mini-setting">Setting1</div>
                        <div className = "mini-setting">Setting2</div>
                    </div>
                </div>
                
            </div>
            
        </>
    );    

}    
Page2.defaultProps = 
{
    Title1: 'App Settings', 
    Title2: 'Shopify Settings',
    subTitle1: 'Sub Title of setting',
    Description1: 'Description of product goes here, as well as any additional information',
    Value1: 'Value of setting currently in the api',

}
export default Page2;