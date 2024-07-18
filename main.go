package main

import (
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"
	pb "github.com/Omkar-Waingankar/serialization/proto" // Update this path as necessary
)

// Sample raw MIME headers
var rawHeaders = map[string]string{
	"Delivered-To":                  "omkar.waingankar@gmail.com",
	"Received":                      "by 2002:a05:7412:7790:b0:11b:95c7:523c with SMTP id ex16csp1596762rdb; Thu, 18 Jul 2024 09:10:59 -0700 (PDT)",
	"X-Google-Smtp-Source":          "AGHT+IHPOEsAr/kp5Gbclu6uZTunWok1qpXAFA16VRnRS7S6FwWeRlAyqEFwgoMPuq4epCZu0FFx",
	"X-Received":                    "by 2002:a05:6214:202c:b0:6b5:238:2e42 with SMTP id 6a1803df08f44-6b79c806f41mr53302786d6.40.1721319059125; Thu, 18 Jul 2024 09:10:59 -0700 (PDT)",
	"ARC-Seal":                      "i=1; a=rsa-sha256; t=1721319059; cv=none; d=google.com; s=arc-20160816; b=DbsyIblPs5wIRzA6hnaooefBWq8yq8r1vwbrUv0MI7LEQzrkUXE8XaFcNzCYcOxJWE +63cX0QcQ0H9ZW2g4HdxufxMQnG13B6hvr0zZ/asYElWmHNvY2y4WG4bIJdifFFkP1fr TptR97Pv3RuDeAnNjaw3GOsv0MggYSbFOPliNOvYbOqQdMIG/59pi+FrgMlKlViiYgr6 4Ox6Yc/arrHrEjk/MsBfzIUIMHMDqzLAEvZ/dZBmx8+029tYBNyma3sA9AW/QobLyXuv A5ryC1XfQCWurzgmthM2oYKwTH3WU/2Fc7d+sxb47e+f3LZ5Mfxy+Y/iw71t0tabVieu 6YVw==",
	"ARC-Message-Signature":         "i=1; a=rsa-sha256; c=relaxed/relaxed; d=google.com; s=arc-20160816; h=to:list-unsubscribe-post:list-unsubscribe:subject:message-id :mime-version:from:date:dkim-signature:dkim-signature; bh=RnYkFnPH1wUNSZwRisokAzVFJFDkre2x39CwD5lGtkc=; fh=5Y2HtxEmyqCVQoCIUDYcxVum5fVuKhemclPKu8Ozm60=; b=bWYERI0OZ1l8z2EJ+JKbZa8vUBbTYQla+ogML+Q7n+IdQ4rZZiDt0crwDCRsLBotAa 0+Ai408G/wF9CBDRdi4ZzkpUq1rhyimDNt0zuywLWDDpkWODTWWUAusJZxSnBloiLSNI ADegWd1Ko4TqT0MxLyuTJXwqxD0rBTd1ssQ0BRNvQLw5MxSgq9nfV+t6gG7UDqT7iTyp tq499NKFZ7NyPOZ3jtodKkmssgtyK0FSng/2UjotIfQby7tegOgOWiu3tpqiANvaBPOq 6R36z1zirvoV46s5ZIeVzK0Y/E3Cz+HX0fEwZylkiaiddUzqNW/j5itlPuJhD5UnvKDG SrLg==; dara=google.com",
	"ARC-Authentication-Results":    "i=1; mx.google.com; dkim=pass header.i=@spotify.com header.s=s1 header.b=YJFhvbg+; dkim=pass header.i=@sendgrid.info header.s=smtpapi header.b=FC9e3MIL; spf=pass (google.com: domain of bounces+54769-d396-omkar.waingankar=gmail.com@em.spotify.com designates 149.72.133.204 as permitted sender) smtp.mailfrom=\"bounces+54769-d396-omkar.waingankar=gmail.com@em.spotify.com\"; dmarc=pass (p=REJECT sp=REJECT dis=NONE) header.from=spotify.com",
	"Return-Path":                   "<bounces+54769-d396-omkar.waingankar=gmail.com@em.spotify.com>",
	"Received-SPF":                  "pass (google.com: domain of bounces+54769-d396-omkar.waingankar=gmail.com@em.spotify.com designates 149.72.133.204 as permitted sender) client-ip=149.72.133.204;",
	"Authentication-Results":        "mx.google.com; dkim=pass header.i=@spotify.com header.s=s1 header.b=YJFhvbg+; dkim=pass header.i=@sendgrid.info header.s=smtpapi header.b=FC9e3MIL; spf=pass (google.com: domain of bounces+54769-d396-omkar.waingankar=gmail.com@em.spotify.com designates 149.72.133.204 as permitted sender) smtp.mailfrom=\"bounces+54769-d396-omkar.waingankar=gmail.com@em.spotify.com\"; dmarc=pass (p=REJECT sp=REJECT dis=NONE) header.from=spotify.com",
	"DKIM-Signature":                "v=1; a=rsa-sha256; c=relaxed/relaxed; d=spotify.com; h=content-type:from:mime-version:subject:list-unsubscribe: list-unsubscribe-post:x-feedback-id:to:cc:content-type:from:subject:to; s=s1; bh=RnYkFnPH1wUNSZwRisokAzVFJFDkre2x39CwD5lGtkc=; b=YJFhvbg+cwsJtO8hqTQ8eNRzQjTBNnNwLqiSiyxYTqV+vzolqYV2Nww6dvGfzh/1tdGI G9WQMIAD2EkGe9xIp6bmxAOxLQyvVv7VMuawccSTDrj0Rfp6Z2aUvJOG/h4pZd7qBNUkQf sKQngPxMWQlVVUQ3WxnjAU1WBCchZrdP5OwYBUCLnVB7WWgDoApGtAV31NosN4fukdraB2 UM7UOQK30xCFtJGdx9exZ2plWfm+75YSt3YTBFNqRRPZPsAO04L87Q4Khgghy6VTUnv61N AM9YFqjTBzYo4OglbcagJXMr6fZ97aHHlIsiW6k4jKu3lIvPni0L7eSZab3dphWg==",
	"Content-Type":                  "multipart/alternative; boundary=6d0cf8719d2ff036ac1aeff7f9423df51b409cabc8d670c1404de4d7fc3e",
	"Date":                          "Thu, 18 Jul 2024 16:10:56 +0000 (UTC)",
	"From":                          "Spotify <no-reply@spotify.com>",
	"Mime-Version":                  "1.0",
	"Message-ID":                    "<NJqE7Bf-RnSPZLtjpkqT9w@geopod-ismtpd-12>",
	"Subject":                       "A special thank you from Don Toliver",
	"List-Unsubscribe-Post":         "List-Unsubscribe=One-Click",
	"X-Feedback-ID":                 "54769:SG",
	"X-SG-EID":                      "u001.ZD5aH33R6V7weoiTHIdZyDuP+A6pZLQL0nRW+xTj31Msb5xf5uAbZIKpGToelTWVR8rVonvMrj0NERWIakIdZIQNWG3HKUGZYRm+qwx1Cf6plmGWnV4LTK9Qns8TKgwXHdRYTx4a+zJ1Z+Ko8c5oEF16JmJ/lfjmOxRHOtReYaWbEj7yFVuD/oqzKYSUEKkq+3xYgF1xJd64OENFuvRtFJ+Z/edvk9vZrdF23jPuvgbu7Spd7rmB0QMoAXK1aIFc",
	"X-SG-ID":                       "u001.SdBcvi+Evd/bQef8eZF3BnMMlgQXm6RhSPDYSp38RDA6P9dNCpEwbb6SESIc221/rUSV03ugX3c64sWeKSTylXRa65b19ow6vRHeDGQMfAlx+Ik3KrHSX24gHkK61ync4rITlIOn3E08SNQHI17d6gpheL8+nbVqBz1NBNZp4+LdyKZ4ksHE1ZmbOI68VghMAYHf+EbQ8/MDsFNeKeu5pg1vkB64UfOPKGosb6TysJf335Dkc7ouuANEUchryQTr62ZKLBtNvgJIaESgmaOHBAE0uPB++LaQfKXYbhd8Tru6XDlhJ4MXnHxuktQSwcNlIDHK0VDFljzZTgdb1+LG3b5WYVTum9Y12jG46qGvzLDDq9ZRfGUc+qMMZRtwfY6J",
	"To":                            "omkar.waingankar@gmail.com",
	"X-Entity-ID":                   "u001.fajWXQMUTuQmJ7EuMgN2yg==",
}

func main() {
	// Serialize using protobuf
	headersProto := &pb.MIMEHeaders{
		Headers: rawHeaders,
	}
	protoData, err := proto.Marshal(headersProto)
	if err != nil {
		log.Fatalf("Failed to marshal protobuf: %v", err)
	}

	// Serialize using JSON
	jsonData, err := json.Marshal(rawHeaders)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Compare the lengths
	protoLength := len(protoData)
	jsonLength := len(jsonData)

	fmt.Printf("Protobuf serialized length: %d bytes\n", protoLength)
	fmt.Printf("JSON serialized length: %d bytes\n", jsonLength)
}

