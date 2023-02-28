const NotificationProvider = require("./notification-provider");
const axios = require("axios");
const FormData = require("form-data");

class Lark extends NotificationProvider {

    name = "lark";

    async send(notification, msg, monitorJSON = null, heartbeatJSON = null) {
        let okMsg = "Sent Successfully.";

        try {
            // let data = {
            //     heartbeat: heartbeatJSON,
            //     monitor: monitorJSON,
            //     msg,
            // };
            let text = {
                text: msg
            }
            let data = {
                msg_type: "text",
                content:text               
            }
            let finalData;
            let config = {};

            if (notification.larkContentType === "form-data") {
                finalData = new FormData();
                finalData.append("data", JSON.stringify(data));

                config = {
                    headers: finalData.getHeaders(),
                };

            } else {
                finalData = data;
            }

            await axios.post(notification.larkURL, finalData, config);
            return okMsg;

        } catch (error) {
            this.throwGeneralAxiosError(error);
        }

    }

}

module.exports = Lark;
