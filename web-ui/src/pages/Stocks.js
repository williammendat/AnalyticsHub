import { useState, useEffect } from "react";
import axios from 'axios';
import "../assets/styles/card.css";
import { useHistory } from "react-router-dom";

function Stocks() {

    const [data, setData] = useState([])
    const history = useHistory();

    const getSymbolPage = (symbol) => {
      let path = `/stocks/${symbol}`; 
      history.push(path);
    }

    useEffect(() => {
        const fetchStocks = async () =>{
          try {
            const {data: response} = await axios.get('http://localhost:2345/stock');
            const models = response.models;
            setData(models);
          } catch (error) {
            console.error(error.message);
          }
        }
    
        fetchStocks();
      }, []);

    const getCards = () => {
        let content = [];
        for (let item of data) {
          content.push
          (
            <div key={item.Id} className="card" onClick={() => getSymbolPage(item.symbol)}>
                <div className="card__image-container">
                    <img
                        src={item.logoURL}
                    />
                </div>

                <div className="card__content">
                    <p className="card__title text--medium">
                        {item.name}
                    </p>
                    <div className="card__info">
                        <p className="text--medium">Favorite</p>
                    </div>
                </div>
            </div>
          );
        }
        return content;
      };

    return (
      <>
        <div className="card-container">
        <section className="cards">
            {getCards()}
        </section>
        </div>
      </>
    );
  }
  
  export default Stocks;
  