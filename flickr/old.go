package flickr

// func (request *Request) buildPost(url_ string, filename string, filetype string) (*http.Request, error) {
// 	real_url, _ := url.Parse(url_)

// 	f, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}

// 	stat, err := f.Stat()
// 	if err != nil {
// 		return nil, err
// 	}
// 	f_size := stat.Size()

// 	request.Args["api_key"] = request.ApiKey

// 	boundary, end := "----###---###--flickr-go-rules", "\r\n"

// 	// Build out all of POST body sans file
// 	header := bytes.NewBuffer(nil)
// 	for k, v := range request.Args {
// 		header.WriteString("--" + boundary + end)
// 		header.WriteString("Content-Disposition: form-data; name=\"" + k + "\"" + end + end)
// 		header.WriteString(v + end)
// 	}
// 	header.WriteString("--" + boundary + end)
// 	header.WriteString("Content-Disposition: form-data; name=\"photo\"; filename=\"photo.jpg\"" + end)
// 	header.WriteString("Content-Type: " + filetype + end + end)

// 	footer := bytes.NewBufferString(end + "--" + boundary + "--" + end)

// 	body_len := int64(header.Len()) + int64(footer.Len()) + f_size

// 	r, w := io.Pipe()
// 	go func() {
// 		pieces := []io.Reader{header, f, footer}

// 		for _, k := range pieces {
// 			_, err = io.Copy(w, k)
// 			if err != nil {
// 				w.CloseWithError(nil)
// 				return
// 			}
// 		}
// 		f.Close()
// 		w.Close()
// 	}()

// 	http_header := make(http.Header)
// 	http_header.Add("Content-Type", "multipart/form-data; boundary="+boundary)

// 	postRequest := &http.Request{
// 		Method:        "POST",
// 		URL:           real_url,
// 		Host:          apiHost,
// 		Header:        http_header,
// 		Body:          r,
// 		ContentLength: body_len,
// 	}
// 	return postRequest, nil
// }

// Example:
// r.Upload("thumb.jpg", "image/jpeg")
// func (request *Request) Upload(filename string, filetype string) (response *Response, err error) {
// 	postRequest, err := request.buildPost(uploadEndpoint, filename, filetype)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return sendPost(postRequest)
// }

// func (request *Request) Replace(filename string, filetype string) (response *Response, err error) {
// 	postRequest, err := request.buildPost(replaceEndpoint, filename, filetype)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return sendPost(postRequest)
// }

// func sendPost(req *http.Request) (response *Response, err error) {
// 	// Create and use TCP connection (lifted mostly wholesale from http.send)
// 	resp, err := http.DefaultClient.Do(req)

// 	if err != nil {
// 		return nil, err
// 	}
// 	rawBody, _ := ioutil.ReadAll(resp.Body)
// 	resp.Body.Close()

// 	var r Response
// 	err = xml.Unmarshal(rawBody, &r)

// 	return &r, err
// }
