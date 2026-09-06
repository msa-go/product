package middleware

import (
	"context"
	"product/infrastructure/log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// http.Request.Context()는 불변이라 기존 요청의 컨텍스트를 직접 바꿀 수 없음.
		// WithContext는 원본은 그대로 두고 새 컨텍스트를 가진 얕은 복사본을 반환하므로,
		// 그 복사본을 c.Request에 재할당해야 이후 핸들러들이 새 컨텍스트를 사용할 수 있음.
		// =>  timeoutCtx(2초 타임아웃)와 request_id를 담아 재할당함
		//
		//  - 타임아웃(timeoutCtx):
		//		이후 핸들러가 c.Request.Context()를 그대로 넘겨 DB 조회/외부 API 호출을 하면, 일정 시간(ex. 2초)이 지났을 때
		//		자동으로 취소(cancel)시켜 느린 하위 요청이 전체 요청을 무한정 붙잡지 않도록 강제할 수 있음
		//
		//  - request_id:
		//		컨트롤러/서비스/리포지토리 등 호출 체인 깊은 곳에서도 ctx.Value("request_id")로 동일한 ID를 꺼내 로그에 남길 수 있어,
		//		하나의 요청에서 발생한 여러 로그를 request_id 기준으로 추적(correlation)할 수 있음
		requestID := uuid.New().String()

		// c.Request.Context()를 parent로 사용: 클라이언트가 연결을 끊으면 즉시 취소되고,
		// 연결이 살아있어도 2초가 지나면 fast-fail로 취소됨 (둘 중 먼저 오는 조건이 적용)
		timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		ctx := context.WithValue(timeoutCtx, "request_id", requestID)

		c.Request = c.Request.WithContext(ctx)

		startTime := time.Now()          // 실행되는 시점(요청이 들어온 직후)의 시간 기록
		c.Next()                         // 제어권을 다음으로 넘김. 여기서 "다음"은 라우터 하나만이 아니라, 체인에 등록된 나머지 미들웨어들 + 최종적으로 매칭된 실제 핸들러(컨트롤러) 전부를 순서대로 실행
		latency := time.Since(startTime) // c.Next()가 반환된 시점(핸들러 및 이후 모든 체인이 끝난 시점)까지의 소요 시간 계산

		requestLog := logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"latency":    latency,
		}

		if c.Writer.Status() < 400 {
			log.Logger.WithFields(requestLog).Info("Request Success")
		} else {
			log.Logger.WithFields(requestLog).Error("Request Failed")
		}
	}
}
