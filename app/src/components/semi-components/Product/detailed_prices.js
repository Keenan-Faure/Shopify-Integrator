import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_Price(props)
{
    return (
        <>
            <div className = "variant-price">{props.Price}</div>
        </>
    );
}

Detailed_Price.defaultProps = 
{
    Price: 'Pricesss'
}
export default Detailed_Price;
