import '../../CSS/detailed.css';

function Detailed_User(props)
{
    return (
        <>
            <table className = "table" style ={{marginBottom: '0px'}}>
                <tbody>
                    <tr>
                        <th style= {{width: '25%'}}>Id</th>
                        <td style= {{width: '25%'}}>{props.id}</td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Name</td>
                        <td style= {{width: '25%'}}>{props.name}</td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Email</td>
                        <td style= {{width: '25%'}}>{props.email}</td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Api_Key</td>
                        <td style= {{width: '25%'}}>{props.name}</td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Webhook_Token</td>
                        <td style= {{width: '25%'}}>{props.web_token}</td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Created At:</td>
                        <td style= {{width: '25%'}}>{props.created}</td>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}}>Updated At:</td>
                        <td style= {{width: '25%'}}>{props.updated}</td>
                    </tr>
                </tbody>
            </table>
        </>
    );
};

export default Detailed_User;

/*
<table className = "table" style ={{marginBottom: '0px'}}>
                <tbody>
                    <tr>
                        <th style= {{width: '25%'}}>Id</th>
                        <th style= {{width: '25%'}}>Name</th>
                        <th style= {{width: '25%'}}>Email</th>
                        <th style= {{width: '25%'}}>Api-Key</th>
                    </tr>
                    <tr>
                        <td style= {{width: '25%'}} className = "price_name">{props.id}</td>
                        <td style= {{width: '25%'}}className = "price_value">{props.name}</td>
                        <td style= {{width: '25%'}} className = "price_name">{props.email}</td>
                        <td style= {{width: '25%'}} className = "price_value">{props.api_key}</td>
                    </tr>
                </tbody>
            </table>

            <table className = "table" style ={{marginBottom: '10px'}}>
                <tbody>
                    <tr>
                        <th style= {{width: '33%'}}>WebHook Token</th>
                        <th style= {{width: '33%'}}>Created At:</th>
                        <th style= {{width: '33%'}}>Updated At:</th>
                    </tr>
                    <tr>
                        <td style= {{width: '33%'}} className = "price_name">{props.web_token}</td>
                        <td style= {{width: '33%'}} className = "price_value">{props.created}</td>
                        <td style= {{width: '33%'}} className = "price_value">{props.updated}</td>
                    </tr>
                </tbody>
            </table>
*/
