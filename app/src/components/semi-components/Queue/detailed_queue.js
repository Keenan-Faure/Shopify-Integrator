import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_queue(props)
{
    useEffect(()=> 
    {
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

        document.getElementById("Customer").click();

        /* When the user clicks on the return button */
        let close = document.querySelector(".rtn-button");
        let filter = document.querySelector(".filter");
        let main = document.querySelector(".main");
        let navbar = document.getElementById("navbar");
        let details = document.querySelector(".details");
        close.addEventListener("click", ()=> 
        {
            close.style.display = "none";
            details.style.animation = "Fadeout 0.5s ease-out";
            
            setTimeout(() => 
            {
                main.style.animation = "FadeIn ease-in 0.6s";
                filter.style.animation = "FadeIn ease-in 0.6s";
                navbar.style.animation = "FadeIn ease-in 0.6s";
                details.style.display = "none";
                navbar.style.display = "block";
                main.style.display = "block";
                filter.style.display = "block";
            }, 500);
        });

    }, []);

    return (
        
        <div id = "detailss" style = {{display: props.Display}}>
            <div className = 'rtn-button'></div>
            <div className = "button-holder">
                <button className="tablink" style={{width: '9%'}} id = "Customer">Queue Item</button>
            </div>
        
            <div className="tabcontent" id="_Customer" >
                <div className = "details-details">
                    <div className = "detailed-image" />
                    <div className = "detailed">
                        <div className = "details-title">{props.Queue_ID}</div>
                        <table>
                            <tbody>
                                <tr>
                                    <th style={{width: '25%'}}>Queue Instruction</th>
                                    <th style={{width: '25%'}}>Queue Type</th>
                                    <th style={{width: '25%'}}>Queue Description</th>
                                    <th style={{width: '25%'}}>Queue Status</th>
                                </tr>
                                <tr>
                                    <td style={{width: '25%'}}>{props.Queue_Instruction}</td>
                                    <td style={{width: '25%'}}>{props.Queue_Type}</td>
                                    <td style={{width: '25%'}}>{props.Queue_Description}</td>
                                    <td style={{width: '25%'}}>{props.Queue_Status}</td>
                                </tr>
                            </tbody>
                        </table> 
                        <div style={{padding: '10px'}}>
                            <div className = "details-description">Created At:</div> 
                            <div style={{marginBottom: '10px'}}>{props.Created_At}</div>
                            <div className = "details-description">Updated At:</div> 
                            <div style={{marginBottom: '10px'}}>{props.Updated_At}</div> 
                        </div>
                        
                    </div>
                    <div className = "details-right"></div>
                    <div className = "details-left"></div>
                </div>
            </div>
        </div>
    );
};

Detailed_queue.defaultProps = 
{
    Queue_Instruction: 'N/A', 
    Queue_Type: 'N/A',
    Queue_Description: 'N/A',
    Queue_Status: 'N/A',
    Queue_ID: 'N/A',

}
export default Detailed_queue;