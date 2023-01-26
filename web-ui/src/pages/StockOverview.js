import { useState, useEffect } from "react";
import axios from "axios";
import { useParams } from "react-router-dom";
import "../assets/styles/StockOverview.css";
import { Chart } from "react-google-charts";
import { LoremIpsum } from "react-lorem-ipsum";

function StockOverview() {

  const [data, setData] = useState(null)
  const [hist, setHist] = useState(null)

  let { id } = useParams();

  useEffect(() => {
    const fetchStocks = async () =>{
      try {
        const url = 'http://localhost:2345/stock/info?symbol=' + id
        console.log(url)
        const {data: response} = await axios.get(url);
        console.log(response.accuracy)
      } catch (error) {
        console.error(error.message);
      }
    }

    fetchStocks();
  }, []);
  const data2 = [
    ["Year", "Sales", "Expenses"],
    ["2004", 1000, 400],
    ["2005", 1170, 460],
    ["2006", 660, 1120],
    ["2007", 1030, 540],
  ];

  const options = {
    title: "Company Performance",
    curveType: "function",
    legend: { position: "bottom" },
  };

  useEffect(() => {
    console.log(id);
  }, []);

  return (
    <>
      <div className="parent">
        <div className="div1">
          <h1>Apple</h1>
          <div>Current open: 150 - Previous open: 160</div>
          <div>- 14%</div>
        </div>
        <div className="div2">
          <LoremIpsum p={2} />
        </div>
        <div className="div3">
          <Chart
            chartType="LineChart"
            width="100%"
            height="100%"
            data={data2}
            options={options}
          />
        </div>
        <div className="div4"> <LoremIpsum p={2} /></div>
        <div className="div5"> <LoremIpsum p={2} /></div>
        <div className="div6"> <LoremIpsum p={2} /></div>
      </div>
    </>
  );
}

export default StockOverview;
